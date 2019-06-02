package stockist

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/beeker1121/goque"
)

var dataDir = "stock_queue"

//Queue to hold stock data before its inserted into the influx db
type Queue struct {
	Client *goque.PrefixQueue
	Prefix string
	Data   map[string]interface{}
	Path   string
}

//New creates a new instance of the Queue
func New(client *goque.PrefixQueue, prefix string, data map[string]interface{}, path string) *Queue {
	queue := &Queue{
		Client: client,
		Prefix: prefix,
		Data:   data,
		Path:   path,
	}

	return queue
}

// Create cretes a new Go Queue
func (q *Queue) Create() (*goque.PrefixQueue, error) {
	c, err := goque.OpenPrefixQueue(q.Path)
	if err != nil {
		fmt.Errorf("Error on creating queue - %v", err)
		return nil, err
	}
	return c, nil

}

//Insert inserts item into the queue
func (q *Queue) Insert() error {
	defer q.Client.Close()
	tickBytes, _ := encode(q.Data)
	_, err := q.Client.Enqueue([]byte(q.Prefix), tickBytes)
	if err != nil {
		fmt.Errorf("Error on inserting into queue - %v", err)
		return err
	}
	fmt.Println("Insert successfull")
	return nil
}

//Pop Pops an item from the queue
func (q *Queue) Pop() (map[string]interface{}, error) {
	defer q.Client.Close()
	ticks, err := q.Client.Dequeue([]byte(q.Prefix))
	if err != nil {
		fmt.Errorf("Error on popping from queue - %v", err)
		return nil, err
	}
	tickData, _ := decode(ticks.Value)
	return tickData, nil

}

//Peek peeks next item from the queue
func (q *Queue) Peek() (map[string]interface{}, error) {
	defer q.Client.Close()
	ticks, err := q.Client.Peek([]byte(q.Prefix))
	if err != nil {
		fmt.Errorf("Error on peeking from queue - %v", err)
		return nil, err
	}
	tickData, _ := decode(ticks.Value)
	return tickData, nil
}

//Drop deletes a queue
func (q *Queue) Drop() error {
	defer q.Client.Close()
	err := q.Client.Drop()
	if err != nil {
		fmt.Errorf("Error on deleting queue - %v", err)
		return err
	}
	return nil
}

func encode(data map[string]interface{}) ([]byte, error) {
	// var b bytes.Buffer
	// e := gob.NewEncoder(&b)
	// if err := e.Encode(data); err != nil {
	// 	panic(err)
	// }
	// return b
	gob.Register(map[string]interface{}{})
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

func decode(b []byte) (map[string]interface{}, error) {
	var tickDecode map[string]interface{}
	buf := bytes.NewBuffer(b)
	d := gob.NewDecoder(buf)
	if err := d.Decode(&tickDecode); err != nil {
		return nil, err
	}

	return tickDecode, nil

}
