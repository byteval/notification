package get

type Request struct {
	ID string `uri:"id" binding:"required,uuid"`
}
