package rabinv1alpha1

import "strconv"

func (d *Deal) GetId() string {
	return strconv.Itoa(int(d.Index))
}

func (r *Response) GetId() string {
	return strconv.Itoa(int(r.Index))
}
