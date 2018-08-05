package permission

import "testing"

func Test_tokenHandler(t *testing.T) {
	p := initService()

	got := p.TokenHasPermission("t5", "permission1")
	want := false
	if got != want {
		t.Errorf("Service.TokenHasPermission() = %v, want %v", got, want)
	}

	tokenHandler(p, []byte("{ \"userId\": \"u1\", \"accessToken\": \"t5\", \"accessTokenExpiresAt\":\"2906-01-02T15:04:05.999999999Z\" }"))

	got = p.TokenHasPermission("t5", "permission1")
	want = true
	if got != want {
		t.Errorf("Service.TokenHasPermission() = %v, want %v", got, want)
	}
}

func Test_roleHandler(t *testing.T) {
	p := initService()

	got := p.TokenHasPermission("t1", "permission5")
	want := false
	if got != want {
		t.Errorf("Service.TokenHasPermission() = %v, want %v", got, want)
	}

	roleHandler(p, []byte("{ \"_id\": \"r1\", \"permissions\": [\"p1\", \"p2\", \"p5\"] }"))

	got = p.TokenHasPermission("t1", "permission5")
	want = true
	if got != want {
		t.Errorf("Service.TokenHasPermission() = %v, want %v", got, want)
	}

	got = p.TokenHasPermission("t1", "permission3")
	want = false
	if got != want {
		t.Errorf("Service.TokenHasPermission() = %v, want %v", got, want)
	}
}
