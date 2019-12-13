package wrapper

type Item struct {
	UUID         string
	TemplateUUID string
	Trashed      string
	CreatedAt    string
	UpdatedAt    string
	ChangerUUID  string
	ItemVersion  int
	VaultUUID    string
	Details      struct {
		Fields []struct {
			Designation string
			ID          string
			Name        string
			Type        string
			Value       string
		}
		HTMLForm struct {
			HTMLMethod string
		}
		Password        string
		PasswordHistory []struct {
			Time  int
			Value string
		}
		Sections []struct {
			Name  string
			Title string
		}
	}
	Overview struct {
		URLs []struct {
			L string
			U string
		}
		Ainfo string
		PBE   float64
		Pgrng bool
		Ps    int
		Tags  []string
		Title string
		URL   string
	}
}

func (i *Item) Username() string {
	return i.Overview.Ainfo
}

func (i *Item) Password() string {
	for _, f := range i.Details.Fields {
		if f.Type == "P" {
			return f.Value
		}
	}
	return ""
}
