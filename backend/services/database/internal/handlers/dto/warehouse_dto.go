package dto

//Input DTO's
type CreateWarehouseRequest struct {
    Name        string      `json:"name"        binding:"required"`
    CompanyID   uint        `json:"company_id"  binding:"required"`
}

type UpdateWarehouseRequest struct {
    Name        *string     `json:"name,omitempty" binding:"omitempty"`
    CompanyID   *uint       `json:"name,omitempty" binding:"omitempty"`
}

type GetWarehouseByIDRequest struct {
    ID          string      ``
}

type DeleteWarehouseRequest struct {

}

type ListWarehousesByCompanyIDRequest struct {

}

type CountWarehousesByCompanyIDRequest struct {

}

//Output DTO's
