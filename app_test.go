package steamapi

import "testing"

var appIdMap = map[int64]string{
	620:    "Portal 2",
	220200: "Kerbal Space Program",
	422970: "Devil Daggers",
}

func TestAppDetailed(t *testing.T) {
	for i := range appIdMap {
		if ad, err := GetAppDetailed(i, "EE", "english"); err == nil {
			t.Logf("App detailed: %+#v", ad)
		} else {
			t.Error(err)
			t.Fail()
		}
	}
}

func TestAppList(t *testing.T) {
	if al, err := GetAppList(); err == nil {
		t.Logf("Apps count: %d", len(al.AppList))
		for id, name := range appIdMap {
			b := false
			for _, a := range al.AppList {
				if (a.AppId == id && a.Name != name) ||
					(a.AppId != id && a.Name == name) {
					t.Errorf("%v dont match '%s':%d case", a, name, id)
				} else if a.AppId == id && a.Name == name {
					b = true
					t.Logf("App %v match '%s':%d case", a, name, id)
				}
			}
			if !b {
				t.Errorf("Cannot find game '%s':%d", name, id)
			}
		}
	} else {
		t.Error(err)
		t.Fail()
	}
}
