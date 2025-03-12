package cartitems

import (
	"github.com/google/uuid"
	cart "github.com/terraskye/vertical-slice-implementation/cart/events"
	"github.com/terraskye/vertical-slice-implementation/infra"
)

// once it is installed, and loaded with the integration
var _ infra.Readmodel = (*ReadModel)(nil)

type Query struct {
	CartId uuid.UUID
}

func (q Query) ID() []byte { return []byte(q.CartId.String()) }

type ReadModel struct {
	AggregateId uuid.UUID
	TotalPrice  float64
	Items       map[uuid.UUID]*CartItem
}

func (p *ReadModel) OnCartCleared(_ *cart.CartCleared) {
	p.TotalPrice = 0
	p.Items = make(map[uuid.UUID]*CartItem)
}

func (p *ReadModel) OnItemArchived(ev *cart.ItemArchived) {
	p.TotalPrice -= p.Items[ev.ItemId].Price
	delete(p.Items, ev.ItemId)
}

func (p *ReadModel) OnCartCreated(ev *cart.CartCreated) {
	p.AggregateId = ev.AggregateId
}

func (p *ReadModel) OnItemAdded(ev *cart.ItemAdded) {
	p.Items[ev.ItemId] = &CartItem{
		Description: ev.Description,
		Image:       ev.Image,
		Price:       ev.Price,
		ProductId:   ev.ProductId,
		ItemId:      ev.ItemId,
	}

	p.TotalPrice += ev.Price
}

func (p *ReadModel) OnItemRemoved(ev *cart.ItemRemoved) {
	p.TotalPrice -= p.Items[ev.ItemId].Price
	delete(p.Items, ev.ItemId)
}

type CartItem struct {
	Description string
	Image       string
	Price       float64
	ProductId   uuid.UUID
	ItemId      uuid.UUID
}
