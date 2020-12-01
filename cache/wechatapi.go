package cache

import "time"

// SetCVT 设置component_verify_ticket
func (me *Cache) SetCVT(cvt string) error {
	return me.rc.SetExpire("c_v_ticket", cvt, time.Hour*12)
}

// GetCVT 获取component_verify_ticket
func (me *Cache) GetCVT() string {
	component_verify_ticket, _ := me.rc.GetString("c_v_ticket")
	return component_verify_ticket
}

// SetCAT 设置component_access_token
func (me *Cache) SetCAT(cvt string, expire int64) error {
	return me.rc.SetExpire("c_a_token", cvt, time.Second*time.Duration(expire))
}

// GetCAT 获取component_access_token
func (me *Cache) GetCAT() string {
	component_access_token, _ := me.rc.GetString("c_a_token")
	return component_access_token
}
