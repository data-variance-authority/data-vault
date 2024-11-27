package internal

// Index is a data structure that stores records and allows for fast search
type Index struct {
	Index InvertedIndex
	Meta  MetaIndex
}

// InvertedIndex is a map of attributes to values to record ids
type InvertedIndex map[string]map[string]map[string]bool

// MetaIndex is a map of record ids to attributes
type MetaIndex map[string]map[string]string

// Record is a record in the index
type Record struct {
	Id         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}

// Add adds a record to the index
func (i Index) Add(r Record) {
	// Add record to the index
	for k, v := range r.Attributes {
		if _, ok := i.Index[k]; !ok {
			i.Index[k] = make(map[string]map[string]bool)
		}
		if _, ok := i.Index[k][v]; !ok {
			i.Index[k][v] = make(map[string]bool)
		}
		i.Index[k][v][r.Id] = true
	}

	// Add record to the meta index
	i.Meta[r.Id] = r.Attributes
}

// Remove removes a record from the index
func (i Index) Remove(r Record) {
	for k, v := range r.Attributes {
		delete(i.Index[k][v], r.Id)
		if len(i.Index[k][v]) == 0 {
			delete(i.Index[k], v)
		}
		if len(i.Index[k]) == 0 {
			delete(i.Index, k)
		}
	}

	delete(i.Meta, r.Id)
}

// Get returns a record from the index by id
func (i Index) Get(id string) Record {
	return Record{
		Id:         id,
		Attributes: i.Meta[id],
	}
}

// GetAttributes returns the attributes of a record
func (i Index) GetAttributes(id string) map[string]string {
	return i.Meta[id]
}

// SearchEvery returns a list of records that match all key-value in the query
func (i Index) SearchEvery(query map[string]string) []Record {
	result := make([]Record, 0)
	for attr, value := range query {
		// If the attribute is not in the index, return an empty result
		if _, ok := i.Index[attr]; !ok {
			return result
		}
		// If the value is not in the index, return an empty result
		if _, ok := i.Index[attr][value]; !ok {
			return result
		}
		// For each record id in the index, add the record to the result
		for id := range i.Index[attr][value] {
			result = append(result, Record{
				Id:         id,
				Attributes: i.Meta[id],
			})
		}
	}

	return result
}

// SearchAny returns a list of records that match any key-value in the query
func (i Index) SearchAny(query map[string]string) []Record {
	result := make([]Record, 0)
	for attr, value := range query {
		if _, ok := i.Index[attr]; !ok {
			continue
		}
		if _, ok := i.Index[attr][value]; !ok {
			continue
		}
		for id := range i.Index[attr][value] {
			result = append(result, Record{
				Id:         id,
				Attributes: i.Meta[id],
			})
		}
	}
	return result
}

// SearchAll returns a list of records that match all key in the query
func (i Index) SearchAll(query []string) []Record {
	result := make([]Record, 0)
	for _, attr := range query {
		if _, ok := i.Index[attr]; !ok {
			return result
		}
		for value := range i.Index[attr] {
			for id := range i.Index[attr][value] {
				result = append(result, Record{
					Id:         id,
					Attributes: i.Meta[id],
				})
			}
		}
	}

	return result
}

// NewIndex creates a new InvertedIndex
func NewIndex() Index {
	return Index{
		Index: make(InvertedIndex),
		Meta:  make(MetaIndex),
	}
}
