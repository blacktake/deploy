package forms

//ProjectForm ...
type ProjectForm struct {
	UserName               string `form:"user_name" json:"user_name" binding:"required"`
	TaskName               string `form:"task_name" json:"task_name" binding:"required"`
	EnvironmentType        string `form:"environment_type" json:"environment_type" binding:"required"`
	Group                  string `form:"group" json:"group" binding:"required"`
	Version                string `form:"version" json:"version" binding:"required"`
	EmailList              string `form:"email_list" json:"email_list,omitempty"`
	FunctionalIntroduction string `form:"functional_introduction" json:"functional_introduction,omitempty"`
}

type PlusProjectForm struct {
	UserName               string `form:"user_name" json:"user_name" binding:"required"`
	TaskName               string `form:"task_name" json:"task_name" binding:"required"`
	EnvironmentType        string `form:"environment_type" json:"environment_type" binding:"required"`
	Group                  string `form:"group" json:"group" binding:"required"`
	Version                string `form:"version" json:"version" binding:"required"`
	EmailList              string `form:"email_list" json:"email_list,omitempty"`
	FunctionalIntroduction string `form:"functional_introduction" json:"functional_introduction,omitempty"`
}
