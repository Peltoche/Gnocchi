package contacts

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Peltoche/halium/internal/tools/uuid"
)

var colors = []string{
	"#ffcdb2",
	"#ffb4a2",
	"#e5989b",
	"#b5838d",
	"#6d6875",
}

type Name struct {
	prefix     string
	firstName  string
	middleName string
	surname    string
	suffix     string
}

func (n *Name) Prefix() string {
	return n.prefix
}

func (n *Name) FirstName() string {
	return n.firstName
}

func (n *Name) MiddleName() string {
	return n.middleName
}

func (n *Name) Surname() string {
	return n.surname
}

func (n *Name) Suffix() string {
	return n.suffix
}

func (n *Name) DisplayName() string {
	res := ""

	if n.firstName == "" && n.middleName == "" && n.surname == "" {
		return "(No name)"
	}

	if n.prefix != "" {
		res = n.prefix + " "
	}

	if n.firstName != "" {
		res = res + n.firstName + " "
	}

	if n.middleName != "" {
		res = res + n.middleName + " "
	}

	if n.surname != "" {
		res = res + n.surname + " "
	}

	if n.suffix != "" {
		res = res + n.suffix + ""
	}

	res = res[:len(res)-1]

	return res
}

type Contact struct {
	createdAt time.Time
	name      *Name
	id        uuid.UUID
}

func (c Contact) ID() uuid.UUID { return c.id }

func (c Contact) Name() *Name { return c.name }

func (c Contact) Color() string {
	// Take the user id first two letters and use them as hexa value.
	twoletterhexa := string(c.id)[:2]
	index, _ := strconv.ParseInt(twoletterhexa, 16, 64)

	fmt.Printf("index: %d -> %d\n\n", len(colors), int(index)%len(colors))
	return colors[int(index)%len(colors)]
}

func (c Contact) CreatedAt() time.Time { return c.createdAt }

type CreateCmd struct{}

type EditNameCmd struct {
	Contact    *Contact
	Prefix     string
	FirstName  string
	MiddleName string
	Surname    string
	Suffix     string
}
