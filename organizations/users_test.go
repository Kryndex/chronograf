package organizations_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/influxdata/chronograf"
	"github.com/influxdata/chronograf/organizations"
)

// IgnoreFields is used because ID is created by BoltDB and cannot be predicted reliably
// EquateEmpty is used because we want nil slices, arrays, and maps to be equal to the empty map
var userCmpOptions = cmp.Options{
	cmpopts.IgnoreFields(chronograf.User{}, "ID"),
	cmpopts.EquateEmpty(),
}

func TestUsersStore_Get(t *testing.T) {
	type args struct {
		ctx   context.Context
		usr   *chronograf.User
		orgID string
	}
	tests := []struct {
		name     string
		args     args
		want     *chronograf.User
		wantErr  bool
		addFirst bool
	}{
		{
			name: "Get user with no role in organization",
			args: args{
				ctx: context.Background(),
				usr: &chronograf.User{
					Name:     "billietta",
					Provider: "google",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "The HillBilliettas",
						},
					},
				},
				orgID: "1336",
			},
			wantErr:  true,
			addFirst: true,
		},
		{
			name: "Get user no organization set",
			args: args{
				ctx: context.Background(),
				usr: &chronograf.User{
					Name:     "billietta",
					Provider: "google",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "The HillBilliettas",
						},
					},
				},
			},
			wantErr:  true,
			addFirst: true,
		},
		{
			name: "Get user scoped to an organization",
			args: args{
				ctx: context.Background(),
				usr: &chronograf.User{
					Name:     "billietta",
					Provider: "google",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "The HillBilliettas",
						},
						{
							Organization: "1336",
							Name:         "The BillHilliettos",
						},
					},
				},
				orgID: "1336",
			},
			want: &chronograf.User{
				Name:     "billietta",
				Provider: "google",
				Scheme:   "oauth2",
				Roles: []chronograf.Role{
					{
						Organization: "1336",
						Name:         "The BillHilliettos",
					},
				},
			},
			addFirst: true,
		},
	}
	for _, tt := range tests {
		client, err := NewTestClient()
		if err != nil {
			t.Fatal(err)
		}
		if err := client.Open(context.TODO()); err != nil {
			t.Fatal(err)
		}
		defer client.Close()

		tt.args.ctx = context.WithValue(tt.args.ctx, "organizationID", tt.args.orgID)
		if tt.addFirst {
			tt.args.usr, err = client.UsersStore.Add(tt.args.ctx, tt.args.usr)
			if err != nil {
				t.Fatal(err)
			}
		}
		s := organizations.NewUsersStore(client.UsersStore)
		got, err := s.Get(tt.args.ctx, chronograf.UserQuery{ID: &tt.args.usr.ID})
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. UsersStore.Get() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if diff := cmp.Diff(got, tt.want, userCmpOptions...); diff != "" {
			t.Errorf("%q. UsersStore.Get():\n-got/+want\ndiff %s", tt.name, diff)
		}
	}
}

func TestUsersStore_Add(t *testing.T) {
	type args struct {
		ctx      context.Context
		u        *chronograf.User
		orgID    string
		uInitial *chronograf.User
	}
	tests := []struct {
		name     string
		args     args
		addFirst bool
		want     *chronograf.User
		wantErr  bool
	}{
		{
			name: "Add new user - no org",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "docbrown",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1336",
							Name:         "editor",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Add new user",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "docbrown",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1336",
							Name:         "editor",
						},
					},
				},
				orgID: "1336",
			},
			want: &chronograf.User{
				Name:     "docbrown",
				Provider: "github",
				Scheme:   "oauth2",
				Roles: []chronograf.Role{
					{
						Organization: "1336",
						Name:         "editor",
					},
				},
			},
		},
		{
			name: "Add non-new user without Role",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "docbrown",
					Provider: "github",
					Scheme:   "oauth2",
					Roles:    []chronograf.Role{},
				},
				orgID: "1336",
				uInitial: &chronograf.User{
					Name:     "docbrown",
					Provider: "github",
					Scheme:   "oauth2",
					Roles:    []chronograf.Role{},
				},
			},
			addFirst: true,
			want: &chronograf.User{
				Name:     "docbrown",
				Provider: "github",
				Scheme:   "oauth2",
				Roles:    []chronograf.Role{},
			},
		},
		{
			name: "Add non-new user with Role",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "docbrown",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1336",
							Name:         "admin",
						},
					},
				},
				orgID: "1336",
				uInitial: &chronograf.User{
					Name:     "docbrown",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1337",
							Name:         "editor",
						},
					},
				},
			},
			addFirst: true,
			want: &chronograf.User{
				Name:     "docbrown",
				Provider: "github",
				Scheme:   "oauth2",
				Roles: []chronograf.Role{
					{
						Organization: "1337",
						Name:         "editor",
					},
					{
						Organization: "1336",
						Name:         "admin",
					},
				},
			},
		},
		{
			name: "Has invalid Role: missing Organization",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "henrietta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Name: "editor",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Has invalid Role: missing Name",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "henrietta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1337",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Has invalid Role: missing Role",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "henrietta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles:    []chronograf.Role{},
				},
			},
			wantErr: true,
		},
		{
			name: "Has invalid Organization",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "henrietta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						chronograf.Role{},
					},
				},
				orgID: "1337",
				uInitial: &chronograf.User{
					Name:     "henrietta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1337",
							Name:         "editor",
						},
					},
				},
			},
			addFirst: true,
			wantErr:  true,
		},
		{
			name: "Organization does not match orgID",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "henrietta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "editor",
						},
					},
				},
				orgID: "1337",
				uInitial: &chronograf.User{
					Name:     "henrietta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1337",
							Name:         "editor",
						},
					},
				},
			},
			addFirst: true,
			wantErr:  true,
		},
		{
			name: "Role Name not specified",
			args: args{
				ctx: context.Background(),
				u: &chronograf.User{
					Name:     "henrietta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1337",
						},
					},
				},
				orgID: "1337",
				uInitial: &chronograf.User{
					Name:     "henrietta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1337",
							Name:         "editor",
						},
					},
				},
			},
			addFirst: true,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		client, err := NewTestClient()
		if err != nil {
			t.Fatal(err)
		}
		if err := client.Open(context.TODO()); err != nil {
			t.Fatal(err)
		}
		defer client.Close()

		tt.args.ctx = context.WithValue(tt.args.ctx, "organizationID", tt.args.orgID)
		s := organizations.NewUsersStore(client.UsersStore)

		if tt.addFirst {
			client.UsersStore.Add(tt.args.ctx, tt.args.uInitial)
		}

		got, err := s.Add(tt.args.ctx, tt.args.u)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. UsersStore.Add() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got == nil && tt.want == nil {
			continue
		}
		got, err = client.UsersStore.Get(tt.args.ctx, chronograf.UserQuery{ID: &got.ID})
		if err != nil {
			t.Fatalf("failed to get added user: %v", err)
		}
		if diff := cmp.Diff(got, tt.want, userCmpOptions...); diff != "" {
			t.Errorf("%q. UsersStore.Add():\n-got/+want\ndiff %s", tt.name, diff)
		}
	}
}

func TestUsersStore_Delete(t *testing.T) {
	type args struct {
		ctx   context.Context
		user  *chronograf.User
		orgID string
	}
	tests := []struct {
		name     string
		args     args
		addFirst bool
		wantErr  bool
		wantRaw  *chronograf.User
	}{
		{
			name: "No such user",
			args: args{
				ctx: context.Background(),
				user: &chronograf.User{
					ID: 10,
				},
				orgID: "1336",
			},
			wantErr: true,
		},
		{
			name: "Derlete user",
			args: args{
				ctx: context.Background(),
				user: &chronograf.User{
					Name: "noone",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "The BillHilliettas",
						},
						{
							Organization: "1336",
							Name:         "The HillBilliettas",
						},
					},
				},
				orgID: "1336",
			},
			addFirst: true,
			wantRaw: &chronograf.User{
				Name: "noone",
				Roles: []chronograf.Role{
					{
						Organization: "1338",
						Name:         "The BillHilliettas",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		client, err := NewTestClient()
		if err != nil {
			t.Fatal(err)
		}
		if err := client.Open(context.TODO()); err != nil {
			t.Fatal(err)
		}
		defer client.Close()

		tt.args.ctx = context.WithValue(tt.args.ctx, "organizationID", tt.args.orgID)
		if tt.addFirst {
			tt.args.user, _ = client.UsersStore.Add(tt.args.ctx, tt.args.user)
		}
		s := organizations.NewUsersStore(client.UsersStore)
		if err := s.Delete(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
			t.Errorf("%q. UsersStore.Delete() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if u, err := s.Get(tt.args.ctx, chronograf.UserQuery{ID: &tt.args.user.ID}); err == nil {
			t.Errorf("%q. Expected error retrieving deleted user, got user %v", tt.name, u)
		}
		gotRaw, _ := client.UsersStore.Get(tt.args.ctx, chronograf.UserQuery{ID: &tt.args.user.ID})
		if diff := cmp.Diff(gotRaw, tt.wantRaw, userCmpOptions...); diff != "" {
			t.Errorf("%q. UsersStore.Delete():\n-got/+want\ndiff %s", tt.name, diff)
		}
	}
}

func TestUsersStore_Update(t *testing.T) {
	type args struct {
		ctx   context.Context
		usr   *chronograf.User
		roles []chronograf.Role
		orgID string
	}
	tests := []struct {
		name     string
		args     args
		addFirst bool
		want     *chronograf.User
		wantRaw  *chronograf.User
		wantErr  bool
	}{
		{
			name: "No such user",
			args: args{
				ctx: context.Background(),
				usr: &chronograf.User{
					ID: 10,
				},
				orgID: "1338",
			},
			wantErr: true,
		},
		{
			name: "Update user role",
			args: args{
				ctx: context.Background(),
				usr: &chronograf.User{
					Name:     "bobetta",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "viewer",
						},
						{
							Organization: "1337",
							Name:         "viewer",
						},
					},
				},
				roles: []chronograf.Role{
					{
						Organization: "1338",
						Name:         "editor",
					},
				},
				orgID: "1338",
			},
			want: &chronograf.User{
				Name:     "bobetta",
				Provider: "github",
				Scheme:   "oauth2",
				Roles: []chronograf.Role{
					{
						Organization: "1338",
						Name:         "editor",
					},
				},
			},
			wantRaw: &chronograf.User{
				Name:     "bobetta",
				Provider: "github",
				Scheme:   "oauth2",
				Roles: []chronograf.Role{
					{
						Organization: "1337",
						Name:         "viewer",
					},
					{
						Organization: "1338",
						Name:         "editor",
					},
				},
			},
			addFirst: true,
		},
	}
	for _, tt := range tests {
		client, err := NewTestClient()
		if err != nil {
			t.Fatal(err)
		}
		if err := client.Open(context.TODO()); err != nil {
			t.Fatal(err)
		}
		defer client.Close()

		tt.args.ctx = context.WithValue(tt.args.ctx, "organizationID", tt.args.orgID)
		if tt.addFirst {
			tt.args.usr, err = client.UsersStore.Add(tt.args.ctx, tt.args.usr)
			if err != nil {
				t.Fatal(err)
			}
		}
		s := organizations.NewUsersStore(client.UsersStore)

		if tt.args.roles != nil {
			tt.args.usr.Roles = tt.args.roles
		}

		if err := s.Update(tt.args.ctx, tt.args.usr); (err != nil) != tt.wantErr {
			t.Errorf("%q. UsersStore.Update() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}

		// for the empty test
		if tt.want == nil {
			continue
		}

		got, err := s.Get(tt.args.ctx, chronograf.UserQuery{ID: &tt.args.usr.ID})
		if err != nil {
			t.Fatalf("failed to get updated user: %v", err)
		}
		if diff := cmp.Diff(got, tt.want, userCmpOptions...); diff != "" {
			t.Errorf("%q. UsersStore.Update():\n-got/+want\ndiff %s", tt.name, diff)
		}
		gotRaw, err := client.UsersStore.Get(tt.args.ctx, chronograf.UserQuery{ID: &tt.args.usr.ID})
		if err != nil {
			t.Fatalf("failed to get updated user: %v", err)
		}
		if diff := cmp.Diff(gotRaw, tt.wantRaw, userCmpOptions...); diff != "" {
			t.Errorf("%q. UsersStore.Update():\n-got/+want\ndiff %s", tt.name, diff)
		}
	}
}

func TestUsersStore_All(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		want     []chronograf.User
		wantRaw  []chronograf.User
		orgID    string
		addFirst bool
		wantErr  bool
	}{
		{
			name: "No users",
			ctx:  context.Background(),
			wantRaw: []chronograf.User{
				{
					Name:     "howdy",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "viewer",
						},
						{
							Organization: "1336",
							Name:         "viewer",
						},
					},
				},
				{
					Name:     "doody2",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1337",
							Name:         "editor",
						},
					},
				},
				{
					Name:     "doody",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "editor",
						},
					},
				},
			},
			orgID: "2330",
		},
		{
			name:  "get all users",
			orgID: "1338",
			ctx:   context.Background(),
			want: []chronograf.User{
				{
					Name:     "howdy",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "viewer",
						},
					},
				},
				{
					Name:     "doody",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "editor",
						},
					},
				},
			},
			wantRaw: []chronograf.User{
				{
					Name:     "howdy",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "viewer",
						},
						{
							Organization: "1336",
							Name:         "viewer",
						},
					},
				},
				{
					Name:     "doody2",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1337",
							Name:         "editor",
						},
					},
				},
				{
					Name:     "doody",
					Provider: "github",
					Scheme:   "oauth2",
					Roles: []chronograf.Role{
						{
							Organization: "1338",
							Name:         "editor",
						},
					},
				},
			},
			addFirst: true,
		},
	}
	for _, tt := range tests {
		client, err := NewTestClient()
		if err != nil {
			t.Fatal(err)
		}
		if err := client.Open(context.TODO()); err != nil {
			t.Fatal(err)
		}
		defer client.Close()

		tt.ctx = context.WithValue(tt.ctx, "organizationID", tt.orgID)
		if tt.addFirst {
			for _, u := range tt.wantRaw {
				client.UsersStore.Add(tt.ctx, &u)
			}
		}
		s := organizations.NewUsersStore(client.UsersStore)
		gots, err := s.All(tt.ctx)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. UsersStore.All() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if diff := cmp.Diff(gots, tt.want, userCmpOptions...); diff != "" {
			t.Errorf("%q. UsersStore.All():\n-got/+want\ndiff %s", tt.name, diff)
		}
	}
}