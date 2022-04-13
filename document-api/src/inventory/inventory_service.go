package inventory

type inventoryService struct {
}

// Service is the interface that provides methods.
type InventoryService interface {
	CreateInventory()
	UpdateInventory()
	InventoryAddProduct()
	DeleteInventory()
	GetInventoryProducts()
	InventoryFilter()
}

func NewInventoryService() InventoryService {
	return &inventoryService{}
}

func (is inventoryService) CreateInventory() {

}
func (is inventoryService) UpdateInventory() {

}
func (is inventoryService) InventoryAddProduct() {

}
func (is inventoryService) DeleteInventory() {

}
func (is inventoryService) GetInventoryProducts() {

}
func (is inventoryService) InventoryFilter() {

}
