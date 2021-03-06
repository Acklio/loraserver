package loraserver

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/brocaar/loraserver/models"
	"github.com/brocaar/lorawan"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
)

// NodeTXPayloadQueueTTL defines the TTL of the node TXPayload queue
var NodeTXPayloadQueueTTL = time.Hour * 24 * 5

const (
	nodeTXPayloadQueueTempl     = "node_tx_queue_%s"
	nodeTXPayloadInProcessTempl = "node_tx_in_process_%s"
)

// createNode creates the given Node.
func createNode(db *sqlx.DB, n models.Node) error {
	_, err := db.Exec("insert into node (dev_eui, app_eui, app_key) values ($1, $2, $3)",
		n.DevEUI[:],
		n.AppEUI[:],
		n.AppKey[:],
	)
	if err == nil {
		log.WithField("dev_eui", n.DevEUI).Info("node created")
	}
	return err
}

// updateNode updates the given Node.
func updateNode(db *sqlx.DB, n models.Node) error {
	res, err := db.Exec("update node set app_eui = $1, app_key = $2, used_dev_nonces = $3 where dev_eui = $4",
		n.AppEUI[:],
		n.AppKey[:],
		n.UsedDevNonces,
		n.DevEUI[:],
	)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return errors.New("DevEUI did not match any rows")
	}
	log.WithField("dev_eui", n.DevEUI).Info("node updated")
	return nil
}

// deleteNode deletes the Node matching the given DevEUI.
func deleteNode(db *sqlx.DB, devEUI lorawan.EUI64) error {
	res, err := db.Exec("delete from node where dev_eui = $1",
		devEUI[:],
	)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return errors.New("DevEUI did not match any rows")
	}
	log.WithField("dev_eui", devEUI).Info("node deleted")
	return nil
}

// getNode returns the Node for the given DevEUI.
func getNode(db *sqlx.DB, devEUI lorawan.EUI64) (models.Node, error) {
	var node models.Node
	return node, db.Get(&node, "select * from node where dev_eui = $1", devEUI[:])
}

// getNodes returns a slice of nodes, sorted by DevEUI.
func getNodes(db *sqlx.DB, limit, offset int) ([]models.Node, error) {
	var nodes []models.Node
	return nodes, db.Select(&nodes, "select * from node order by dev_eui limit $1 offset $2", limit, offset)
}

// addTXPayloadToQueue adds the given TXPayload to the queue.
func addTXPayloadToQueue(p *redis.Pool, payload models.TXPayload) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(payload); err != nil {
		return err
	}

	c := p.Get()
	defer c.Close()

	exp := int64(NodeTXPayloadQueueTTL) / int64(time.Millisecond)
	key := fmt.Sprintf(nodeTXPayloadQueueTempl, payload.DevEUI)

	c.Send("MULTI")
	c.Send("LPUSH", key, buf.Bytes())
	c.Send("PEXPIRE", key, exp)
	_, err := c.Do("EXEC")

	if err != nil {
		return err
	}

	log.WithField("dev_eui", payload.DevEUI).Info("payload added to queue")

	return nil
}

// getTXPayloadAndRemainingFromQueue returns the first TXPayload to send
// to the node and a bool indicating if there are more payloads pending to be
// send. This is either an in-process item (e.g. an item that needs
// to be re-transmitted or an item from the queue (which will be marked as
// in-process when consuming). After a successful transmission, don't forget
// to call clearInProcessTXPayload.
// errDoesNotExist is returned when there are no items to send.
func getTXPayloadAndRemainingFromQueue(p *redis.Pool, devEUI lorawan.EUI64) (models.TXPayload, bool, error) {
	var txPayload models.TXPayload
	queueKey := fmt.Sprintf(nodeTXPayloadQueueTempl, devEUI)
	pendingKey := fmt.Sprintf(nodeTXPayloadInProcessTempl, devEUI)

	c := p.Get()
	defer c.Close()

	c.Send("LINDEX", pendingKey, -1)
	c.Send("LLEN", queueKey)
	c.Flush()

	b, err := redis.Bytes(c.Receive())
	if err != nil && err != redis.ErrNil { // something went wrong
		return txPayload, false, err
	}
	if err == nil { // there is an in-process item
		// read the queue size
		i, err := redis.Int(c.Receive())
		if err != nil {
			return txPayload, false, err
		}
		return txPayload, i > 0, gob.NewDecoder(bytes.NewReader(b)).Decode(&txPayload)
	}

	// redis.ErrNil error was returned, return item from the queue
	i, err := redis.Int(c.Receive())
	if i == 0 {
		return txPayload, false, errDoesNotExist
	}

	c.Send("RPOPLPUSH", queueKey, pendingKey)
	c.Send("LLEN", queueKey)
	c.Flush()

	// read payload from queue
	b, err = redis.Bytes(c.Receive())
	if err != nil {
		if err == redis.ErrNil {
			err = errDoesNotExist
		}
		return txPayload, false, err
	}
	// read remaining items
	i, err = redis.Int(c.Receive())
	if err != nil {
		return txPayload, false, err
	}

	return txPayload, i > 0, gob.NewDecoder(bytes.NewReader(b)).Decode(&txPayload)
}

// clearInProcessTXPayload clears the in-process TXPayload (to be called
// after a successful transmission).
func clearInProcessTXPayload(p *redis.Pool, devEUI lorawan.EUI64) error {
	key := fmt.Sprintf(nodeTXPayloadInProcessTempl, devEUI)
	c := p.Get()
	defer c.Close()
	_, err := redis.Int(c.Do("DEL", key))
	if err != nil {
		return err
	}
	log.WithField("dev_eui", devEUI).Info("in-process payload removed")
	return nil
}

// NodeAPI exports the Node related functions.
type NodeAPI struct {
	ctx Context
}

// NewNodeAPI creates a new NodeAPI.
func NewNodeAPI(ctx Context) *NodeAPI {
	return &NodeAPI{
		ctx: ctx,
	}
}

// Get returns the Node for the given DevEUI.
func (a *NodeAPI) Get(devEUI lorawan.EUI64, node *models.Node) error {
	var err error
	*node, err = getNode(a.ctx.DB, devEUI)
	return err
}

// GetList returns a list of nodes (given a limit and offset).
func (a *NodeAPI) GetList(req models.GetListRequest, nodes *[]models.Node) error {
	var err error
	*nodes, err = getNodes(a.ctx.DB, req.Limit, req.Offset)
	return err
}

// Create creates the given Node.
func (a *NodeAPI) Create(node models.Node, devEUI *lorawan.EUI64) error {
	if err := createNode(a.ctx.DB, node); err != nil {
		return err
	}
	*devEUI = node.DevEUI
	return nil
}

// Update updatest the given Node.
func (a *NodeAPI) Update(node models.Node, devEUI *lorawan.EUI64) error {
	if err := updateNode(a.ctx.DB, node); err != nil {
		return err
	}
	*devEUI = node.DevEUI
	return nil
}

// Delete deletes the node matching the given DevEUI.
func (a *NodeAPI) Delete(devEUI lorawan.EUI64, deletedDevEUI *lorawan.EUI64) error {
	if err := deleteNode(a.ctx.DB, devEUI); err != nil {
		return err
	}
	*deletedDevEUI = devEUI
	return nil
}
