```go

// &x   → get address of x
// *T   → pointer type (pointer to T)
// *p   → dereference pointer (get value)
// nil  → pointer points to nothing
//
// Use pointers when:
// - you need to modify a value inside a function
// - you need to distinguish "not provided" vs "zero value"
// - copying large structs is expensive
//
// Do NOT use pointers when:
// - value is small and immutable
// - nil state is not needed
// - readability becomes worse

//Format
// %s	string
// %d	int
// %f	float
// %v	any value
// %+v	struct 
// %T	type
// %p	pointer


//go run ./cmd/api 
```