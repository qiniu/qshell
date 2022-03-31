package flow

type WorkPackage struct {
	WorkRecords []*WorkRecord
}

func (w *WorkPackage) WorkId() string {
	return ""
}

type WorkPacker struct {
	MaxWorkCountPerPackage int
	workPackage            *WorkPackage
}

func NewWorkPacker(maxWorkCountPerPackage int) *WorkPacker {
	if maxWorkCountPerPackage < 1 {
		maxWorkCountPerPackage = 1
	}

	return &WorkPacker{
		MaxWorkCountPerPackage: maxWorkCountPerPackage,
		workPackage: &WorkPackage{
			WorkRecords: make([]*WorkRecord, 0, maxWorkCountPerPackage),
		},
	}
}

func (w *WorkPacker) Pack(work Work) error {
	w.workPackage.WorkRecords = append(w.workPackage.WorkRecords, &WorkRecord{
		Work:   work,
		Result: nil,
		Err:    nil,
	})
	return nil
}

func (w *WorkPacker) GetWorkPackageAndClean(force bool) (p *WorkPackage) {
	if len(w.workPackage.WorkRecords) < w.MaxWorkCountPerPackage && !force {
		return nil
	}

	p = w.workPackage
	w.workPackage = &WorkPackage{
		WorkRecords: make([]*WorkRecord, 0, w.MaxWorkCountPerPackage),
	}
	return
}
