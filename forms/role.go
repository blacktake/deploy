package forms

//AddRoleForm ...
type AddRoleForm struct {
	Rolename    string `form:"rolename" json:"rolename" binding:"required"`
	Description string `form:"description" json:"description"`
}
