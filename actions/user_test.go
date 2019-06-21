package actions

import (
	"github.com/apimsprdev/dao"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestSignIn(t *testing.T) {
	id := dao.CreateUser("test", "testEm@epsi.fr")
	user := dao.GetUserByEmailAndPassword("testPass", "testEm@epsi.fr")
	if id != user.ID.String() {
		t.Errorf("Id was incorrect, got : %s, want %s.", user.ID, id)
	}
	dao.DeleteUser(user.ID)
}
