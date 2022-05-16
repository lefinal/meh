package meh

// ErrorUnwrapper is used for iterating over the recursive structure of wrapped
// errors in Error.WrappedErr.
type ErrorUnwrapper struct {
	// awaitFirstNext is a flag used for skipping the first unwrap. This is needed
	// because otherwise we would need to wrap the initial error.
	awaitFirstNext bool
	// current is the error for the current level that is going to be unwrapped when
	// calling Next.
	current error
	// currentLevel is a counter for keeping track of the current level. The first
	// Next-call will make the current level 0 and increment afterwards.
	currentLevel int
}

// NewErrorUnwrapper allows iterating the given error from top to bottom level
// using ErrorUnwrapper.Next in a loop and getting the current level's error
// using ErrorUnwrapper.Current.
func NewErrorUnwrapper(err error) *ErrorUnwrapper {
	return &ErrorUnwrapper{
		awaitFirstNext: err != nil,
		current:        err,
		currentLevel:   -1,
	}
}

// Next unwraps the current error.
//
// Warning: You always need to make a Next call before being able to retrieve
// the first error via Current. The reason for this is that it allows simple
// iterating in a for-loop using Next as condition like when scanning rows using
// the sql-package.
func (it *ErrorUnwrapper) Next() bool {
	it.currentLevel++
	if it.awaitFirstNext {
		it.awaitFirstNext = false
		return true
	}
	if e, ok := it.current.(*Error); ok {
		it.current = e.WrappedErr
	} else {
		it.current = nil
	}
	return it.current != nil
}

// Current returns the current error. Remember to call Next before the first
// Current call.
func (it *ErrorUnwrapper) Current() error {
	// Check if the first next call was made and return nil if missed in order to
	// avoid bad practise when using the unwrapper.
	if it.awaitFirstNext {
		return nil
	}
	return it.current
}

// Level returns the current level. Starting at -1, it will increment from 0 for
// each Next-call.
func (it *ErrorUnwrapper) Level() int {
	return it.currentLevel
}
