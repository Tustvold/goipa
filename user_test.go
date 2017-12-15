// Copyright 2015 Andrew E. Bruno. All rights reserved.
// Use of this source code is governed by a BSD style
// license that can be found in the LICENSE file.

package ipa

import (
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	c := newTestClient()
	user := os.Getenv("GOIPA_TEST_USER")
	pass := os.Getenv("GOIPA_TEST_PASSWD")
	sess, err := c.Login(user, pass)
	if err != nil {
		t.Error(err)
	}

	if len(sess) == 0 {
		t.Error(err)
	}
}

func TestCreateDeleteUser(t *testing.T) {
	c := newTestClient()
	user := os.Getenv("GOIPA_ADMIN_USER")
	pass := os.Getenv("GOIPA_ADMIN_PASSWD")
	sess, err := c.Login(user, pass)
	if err != nil {
		t.Error(err)
	}

	if len(sess) == 0 {
		t.Error(err)
	}

	createUid := os.Getenv("GOIPA_CREATE_USER_UID")
	createFirst := os.Getenv("GOIPA_CREATE_USER_FIRST")
	createLast := os.Getenv("GOIPA_CREATE_USER_LAST")

	rec, err := c.CreateUser(createUid, createFirst, createLast)

	if err != nil {
		t.Error(err)
	}

	if string(rec.Uid) != createUid {
		t.Errorf("Invalid username")
	}

	if string(rec.First) != createFirst {
		t.Errorf("Invalid first name")
	}

	if string(rec.Last) != createLast {
		t.Errorf("Invalid last name")
	}

	rec2, err := c.GetUserByUidNumber(string(rec.UidNumber))

	if err != nil {
		t.Error(err)
	}

	if rec2.Uid != rec.Uid {
		t.Errorf("Equality Failure")
	}

	err = c.UserUpdateEmail(createUid, "test@five.ai")

	if err != nil {
		t.Error(err)
	}

	err = c.DeleteUser(createUid)

	if err != nil {
		t.Error(err)
	}
}

func TestUserExists(t *testing.T) {
	c := newTestClient()

	user := os.Getenv("GOIPA_TEST_USER")
	pass := os.Getenv("GOIPA_TEST_PASSWD")
	_, err := c.Login(user, pass)
	if err != nil {
		t.Error(err)
	}

	exists, err := c.UserExists(user)

	if err != nil {
		t.Error(err)
	}

	if !exists {
		t.Error("User should exist")
	}

	exists, err = c.UserExists("safjhkgfdhjkfsehkjfsd")

	if err != nil {
		t.Error(err)
	}

	if exists {
		t.Error("User should not exist")
	}
}

func TestGetUser(t *testing.T) {
	c := newTestClient()

	user := os.Getenv("GOIPA_TEST_USER")
	pass := os.Getenv("GOIPA_TEST_PASSWD")
	_, err := c.Login(user, pass)
	if err != nil {
		t.Error(err)
	}

	// Test using ipa_session
	rec, err := c.GetUser(user)

	if err != nil {
		t.Error(err)
	}

	if string(rec.Uid) != user {
		t.Errorf("Invalid user")
	}

	if len(os.Getenv("GOIPA_TEST_KEYTAB")) > 0 {
		c.ClearSession()

		// Test using keytab if set
		rec, err := c.GetUser(user)

		if err != nil || rec == nil {
			t.Error(err)
		}

		if string(rec.Uid) != user {
			t.Errorf("Invalid user")
		}
	}
}

func TestUpdateSSHPubKeys(t *testing.T) {
	c := newTestClient()

	user := os.Getenv("GOIPA_TEST_USER")
	pass := os.Getenv("GOIPA_TEST_PASSWD")
	_, err := c.Login(user, pass)
	if err != nil {
		t.Error(err)
	}

	// Remove any existing public keys
	fp, err := c.UpdateSSHPubKeys(user, []string{})
	if err != nil {
		t.Errorf("Failed to remove existing ssh public keys: %s", err)
	}

	if len(fp) != 0 {
		t.Error("Invalid number of fingerprints returned")
	}

	_, err = c.UpdateSSHPubKeys(user, []string{"invalid key"})
	if err == nil {
		t.Error("Invalid key was updated")
	}

	fp, err = c.UpdateSSHPubKeys(user, []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDVBSs8RP8KPbdMwOmuKgjScx301k1mBZTubfcJc7HKcJ19f1Z/eJ5y9R7LjhsK1WGn8ISRtP2c0NUNPWcZHdWzTv6m2AFL4qniXr2vvKcewq2fxy8uXnUSvS054wwFDW6trmWV1Vrrab0eXO9S7tGGLdx2ySQ8Bzfe8wY3M2/N1gd5dzGSVg3qFspgikTKjRt5rfaWoN+/OWLDg1HHEWjY0Hgqry1bJW3U83SlIi9+JwKW0zxunwImgFsI1xC15lf7X9LOE9e6XGT1km/NTPOqoAvaCCA0KyAK7P6cLjFVAA/k9UnC/QX6JKXoURFRdhPEdFqauF3Xw9rwDFCFkMUp test@localhost"})
	if err != nil {
		t.Error(err)
	}

	if len(fp) != 1 {
		t.Errorf("Wrong number of fingerprints returned")
	}

	if fp[0] != "SHA256:9NiBLAynn/9d9lNcu/rOh5VXdXIJeA1oJDxfBGsI9xc test@localhost (ssh-rsa)" {
		t.Errorf("Invalid fingerprint: Got %s", fp[0])
	}

	// Remove test public keys
	_, err = c.UpdateSSHPubKeys(user, []string{})
	if err != nil {
		t.Error("Failed to remove testing ssh public keys")
	}
}

func TestUpdateMobile(t *testing.T) {
	c := newTestClient()

	user := os.Getenv("GOIPA_TEST_USER")
	pass := os.Getenv("GOIPA_TEST_PASSWD")
	_, err := c.Login(user, pass)
	if err != nil {
		t.Error(err)
	}

	err = c.UserUpdateMobileNumber(user, "")
	if err != nil {
		t.Error("Failed to remove existing mobile number")
	}

	err = c.UserUpdateMobileNumber(user, "+9999999999")
	if err != nil {
		t.Error(err)
	}

	rec, err := c.GetUser(user)
	if err != nil {
		t.Error(err)
	}

	if string(rec.Mobile) != "+9999999999" {
		t.Errorf("Invalid mobile number")
	}
}

func TestUserAuthTypes(t *testing.T) {
	if len(os.Getenv("GOIPA_TEST_KEYTAB")) > 0 {
		c := newTestClient()

		user := os.Getenv("GOIPA_TEST_USER")

		err := c.SetAuthTypes(user, []string{"otp"})
		if err != nil {
			t.Error(err)
		}

		rec, err := c.GetUser(user)
		if err != nil {
			t.Error(err)
		}

		if !rec.OTPOnly() {
			t.Errorf("User auth type should only be OTP")
		}

		err = c.SetAuthTypes(user, nil)
		if err != nil {
			t.Error("Failed to remove existing auth types")
		}
	}
}
