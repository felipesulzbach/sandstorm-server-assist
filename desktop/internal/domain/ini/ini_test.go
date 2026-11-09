package ini

import "testing"

const sample = "[Rcon]\r\nbEnabled=True\r\nPassword=abc\r\nListenPort=27015\r\n\r\n[/Script/Insurgency.INSMultiplayerMode]\r\nbAllowFriendlyFire=True\r\n"

func TestGetSetRemove(t *testing.T) {
	if v, ok := Get(sample, "rcon", "password"); !ok || v != "abc" {
		t.Fatalf("Get = %q %v", v, ok)
	}
	out, changed := Remove(sample, "Rcon", "Password")
	if !changed {
		t.Fatal("deveria remover")
	}
	if _, ok := Get(out, "Rcon", "Password"); ok {
		t.Fatal("senha continua")
	}
	if v, _ := Get(out, "/Script/Insurgency.INSMultiplayerMode", "bAllowFriendlyFire"); v != "True" {
		t.Fatal("outra seção foi afetada")
	}
	out = Set(out, "Rcon", "bEnabled", "False")
	if v, _ := Get(out, "Rcon", "bEnabled"); v != "False" {
		t.Fatal("Set não trocou")
	}
	out = Set(out, "Nova", "Chave", "1")
	if v, _ := Get(out, "Nova", "Chave"); v != "1" {
		t.Fatal("Set não criou seção")
	}
	out = Set(out, "Rcon", "Extra", "x")
	if v, _ := Get(out, "Rcon", "Extra"); v != "x" {
		t.Fatal("Set não criou chave na seção")
	}
}
