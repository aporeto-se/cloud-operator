package types

// import (
// 	"testing"
// )

// func TestProcessOptions(t *testing.T) {

// 	p := NewProcessOptions().
// 		AddOpCode(OpCodeCreateKubeNamespace).
// 		AddOpCode(OpCodeCreateKubeDNSNet).
// 		AddOpCode(OpCodeCreateKubeDNSPolicy).
// 		AddOpCode(OpCodeCreateKubeNodeExtNet).
// 		AddOpCode(OpCodeInstallKubeEnforcer)

// 	if !p.HasOp(OpCodeCreateKubeNamespace) {
// 		t.Fatalf("%s should exist in ops", OpCodeCreateKubeNamespace)
// 	}

// 	if p.HasOp(OpCodeCreateKubeAPINet) {
// 		t.Fatalf("%s should not exist in ops", OpCodeCreateKubeAPINet)
// 	}

// 	if !p.HasOp(OpCodeCreateKubeDNSNet) {
// 		t.Fatalf("%s should exist in ops", OpCodeCreateKubeDNSNet)
// 	}

// 	if !p.HasOp(OpCodeCreateKubeDNSPolicy) {
// 		t.Fatalf("%s should exist in ops", OpCodeCreateKubeDNSPolicy)
// 	}

// 	if !p.HasOp(OpCodeCreateKubeNodeExtNet) {
// 		t.Fatalf("%s should exist in ops", OpCodeCreateKubeNodeExtNet)
// 	}

// 	if !p.HasOp(OpCodeInstallKubeEnforcer) {
// 		t.Fatalf("%s should exist in ops", OpCodeInstallKubeEnforcer)
// 	}

// 	p.AddOpCode(OpCodeCreateKubeAPINet)

// 	if !p.HasOp(OpCodeCreateKubeNamespace) {
// 		t.Fatalf("%s should exist in ops", OpCodeCreateKubeNamespace)
// 	}

// 	if !p.HasOp(OpCodeCreateKubeAPINet) {
// 		t.Fatalf("%s should exist in ops", OpCodeCreateKubeAPINet)
// 	}

// 	if !p.HasOp(OpCodeCreateKubeDNSNet) {
// 		t.Fatalf("%s should exist in ops", OpCodeCreateKubeDNSNet)
// 	}

// 	if !p.HasOp(OpCodeCreateKubeDNSPolicy) {
// 		t.Fatalf("%s should exist in ops", OpCodeCreateKubeDNSPolicy)
// 	}

// 	if !p.HasOp(OpCodeCreateKubeNodeExtNet) {
// 		t.Fatalf("%s should exist in ops", OpCodeCreateKubeNodeExtNet)
// 	}

// 	if !p.HasOp(OpCodeInstallKubeEnforcer) {
// 		t.Fatalf("%s should exist in ops", OpCodeInstallKubeEnforcer)
// 	}

// }

// func TestConfigVersion(t *testing.T) {

// 	configVersion, err := ConfigVersionFromString("V1")

// 	if err != nil {
// 		t.Fatalf("err should not have been returned")
// 	}

// 	if configVersion != ConfigVersionV1 {
// 		t.Fatalf("configVersion should match")
// 	}

// 	_, err = ConfigVersionFromString("this_is_not_valid")
// 	if err == nil {
// 		t.Fatalf("err should have been returned")
// 	}

// }
