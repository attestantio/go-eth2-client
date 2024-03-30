package ssz

import "fmt"

func (d *DynSsz) getSpecValue(name string) (bool, uint64, error) {

	switch name {
	// there are some calculated values, but adding a parser & dynamic calculations for these seems a bit overkill
	case "SYNC_COMMITTEE_SIZE/8":
		ok, val, err := d.getSpecValue("SYNC_COMMITTEE_SIZE")
		if ok {
			return ok, val / 8, err
		}
	case "SYNC_COMMITTEE_SIZE/SYNC_COMMITTEE_SUBNET_COUNT":
		ok1, val1, err1 := d.getSpecValue("SYNC_COMMITTEE_SIZE")
		ok2, val2, err2 := d.getSpecValue("SYNC_COMMITTEE_SUBNET_COUNT")
		if err1 != nil {
			return false, 0, err1
		}
		if err2 != nil {
			return false, 0, err2
		}
		if ok1 && ok2 {
			return true, val1 / val2, nil
		}
	case "DEPOSIT_CONTRACT_TREE_DEPTH+1":
		ok, val, err := d.getSpecValue("DEPOSIT_CONTRACT_TREE_DEPTH")
		if ok {
			return ok, val + 1, err
		}
	default:
		specVal := d.SpecValues[name]
		if specVal != nil {
			specInt, ok := specVal.(uint64)
			if !ok {
				return false, 0, fmt.Errorf("value is not uint64")
			}
			return true, specInt, nil
		}

	}
	return false, 0, nil
}
