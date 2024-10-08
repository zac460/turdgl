package turdgl

import (
	"fmt"
	"image/color"
)

// hoverable is an interface for shapes can detect cursor hovering.
type hoverable interface {
	Shape
	IsWithin(Vec) bool
}

// Button can be build on top of shapes to create pressable buttons.
type Button struct {
	Shape     hoverable        // the base shape the button is built on
	Label     *Text            // the text to display on the button (if any)
	CB        func(MouseState) // the callback function to execute on press
	Trigger   MouseState       // which mouse button must be used to press the button
	Behaviour ButtonBehaviour  // how the button responds to being pressed

	prevMouseState MouseState
	prevMouseLoc   Vec
	prevLabel      string
}

// NewButton constructs a new button from any shape that satisfies the buttonable interface.
func NewButton(shape hoverable, fontPath string) *Button {
	return &Button{
		Shape:     shape,
		Label:     NewText("", shape.GetPos(), fontPath),
		CB:        func(MouseState) { fmt.Println("Warning: Button callback not configured") },
		Trigger:   LeftClick,
		Behaviour: OnAll,
	}
}

// SetCallback configures a callback function to execute every time a press
// or unpress event occurs. The type of event (left-click, right-click, etc...)
// is passed into the function so the callback can take appropriate action.
func (b *Button) SetCallback(callback func(MouseState)) *Button {
	b.CB = callback
	return b
}

// Draw draws the button onto the frame buffer.
func (b *Button) Draw(buf *FrameBuffer) {
	b.Shape.Draw(buf)

	// Align to centre of underlying shape
	b.Label.SetPos(func() Vec {
		switch b.Shape.(type) {
		case *Rect:
			p := b.Shape.GetPos()
			return Vec{p.X + b.Shape.Width()/2, p.Y + b.Shape.Height()/2}
		default:
			return b.Shape.GetPos()
		}
	}())

	b.Label.Draw(buf)
}

// ButtonBehaviour represents how a button responds to being pressed.
type ButtonBehaviour int

const (
	OnAll             ButtonBehaviour = iota // execute behaviour every time Update() is called
	OnPress                                  // execute behaviour on press
	OnRelease                                // execute behaviour on release
	OnPressAndRelease                        // execute behaviour on press and release
	OnHold                                   // execute behaviour as long as button is held down
	OnHover                                  // execute behaviour if cursor is over button
)

// Update examines button state and executes behaviour accordingly.
func (b *Button) Update(win *Window) {

	currentMouseState := win.MouseButtonState()
	hovering := b.Shape.IsWithin(win.MouseLocation())

	switch b.Behaviour {
	case OnAll:
		b.CB(currentMouseState)
	case OnHover:
		if hovering {
			b.CB(currentMouseState)
		}
	case OnPress:
		if hovering && (b.prevMouseState == NoClick && currentMouseState == b.Trigger) {
			b.CB(currentMouseState)
		}
	case OnRelease:
		if hovering && (b.prevMouseState == b.Trigger && currentMouseState == NoClick) {
			b.CB(currentMouseState)
		}
	case OnPressAndRelease:
		if hovering && (b.prevMouseState != currentMouseState) {
			b.CB(currentMouseState)
		}
	case OnHold:
		if hovering && (currentMouseState == b.Trigger) {
			b.CB(currentMouseState)
		}
	default:
		panic("unsupported button behaviour")
	}

	b.prevMouseState = win.MouseButtonState()
	b.prevMouseLoc = win.MouseLocation()
}

// Move moves the button by a given vector.
func (b *Button) Move(mov Vec) {
	b.Shape.Move(mov)
	b.Label.Move(mov)
}

// IsHovering returns whether the cursor is hovering over the button.
func (b *Button) IsHovering() bool {
	return b.Shape.IsWithin(b.prevMouseLoc)
}

// SetLabelText sets the text label to the given string.
func (b *Button) SetLabelText(s string) *Button {
	b.Label.SetText(s)
	return b
}

// SetLabelAlignment sets the alignment of the text label relative to the centre of the shape.
func (b *Button) SetLabelAlignment(align Alignment) *Button {
	b.Label.SetAlignment(align)
	return b
}

// SetLabelOffset manually sets the label's offset, providing the text is in AlignCustom mode.
func (b *Button) SetLabelOffset(offset Vec) *Button {
	b.Label.SetOffset(offset)
	return b
}

// SetLabelPos sets the label's position on the screen.
func (b *Button) SetLabelPos(pos Vec) *Button {
	b.Label.SetPos(pos)
	return b
}

// SetLabelColour sets the label text colour.
func (b *Button) SetLabelColour(c color.Color) *Button {
	b.Label.SetColour(c)
	return b
}

// SetLabelFont sets the path fo the .ttf file that is used to generate the label.
func (b *Button) SetLabelFont(path string) *Button {
	b.Label.SetFont(path)
	return b
}

// SetLabelDPI sets the DPI of the label font.
func (b *Button) SetLabelDPI(dpi float64) *Button {
	b.Label.SetDPI(dpi)
	return b
}

// SetLabelSize sets the size of the label font.
func (b *Button) SetLabelSize(size float64) *Button {
	b.Label.SetSize(size)
	return b
}

// SetLabelSpacing sets the line spacing of the label.
func (b *Button) SetLabelSpacing(spacing float64) *Button {
	b.Label.SetSpacing(spacing)
	return b
}

// SetLabelMaskSize sets the size of the mask used to generate the label.
func (b *Button) SetLabelMaskSize(w, h int) *Button {
	b.Label.SetMaskSize(w, h)
	return b
}
