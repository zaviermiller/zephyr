package zephyr

import "syscall/js"

type Emitter interface {
	Call()
}

type DOMEvent struct {
	Target js.Value

	Bubbles          bool
	DefaultPrevented bool
}

func eventRecurSet(node *VNode, q chan DOMUpdate) {
	// fmt.Println("yo: ", node.FirstChild, node.events)
	if node.events != nil {
		q <- DOMUpdate{Operation: AddEventListeners, ElementID: node.DOM_ID, Data: node.events}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		eventRecurSet(c, q)
	}
}
func (e *DOMEvent) StopPropagation() {
	e.Bubbles = false
}

func (e *DOMEvent) PreventDefault() {
	e.DefaultPrevented = true
}

var domEvents map[string]struct{} = map[string]struct{}{
	"afterprint":               void,
	"animationend":             void,
	"animationiteration":       void,
	"animationstart":           void,
	"appinstalled":             void,
	"audioprocess":             void,
	"audioend":                 void,
	"audiostart":               void,
	"beforeprint":              void,
	"beforeunload":             void,
	"canplay":                  void,
	"canplaythrough":           void,
	"chargingchange":           void,
	"chargingtimechange":       void,
	"compassneedscalibration":  void,
	"compositionend":           void,
	"compositionstart":         void,
	"compositionupdate":        void,
	"contextmenu":              void,
	"dblclick":                 void,
	"devicechange":             void,
	"devicelight":              void,
	"devicemotion":             void,
	"deviceorientation":        void,
	"deviceproximity":          void,
	"dischargingtimechange":    void,
	"dragend":                  void,
	"dragenter":                void,
	"dragleave":                void,
	"dragover":                 void,
	"dragstart":                void,
	"durationchange":           void,
	"focusin":                  void,
	"focusout":                 void,
	"fullscreenchange":         void,
	"fullscreenerror":          void,
	"gamepadconnected":         void,
	"gamepaddisconnected":      void,
	"gotpointercapture":        void,
	"hashchange":               void,
	"keydown":                  void,
	"keypress":                 void,
	"keyup":                    void,
	"languagechange":           void,
	"levelchange":              void,
	"loadeddata":               void,
	"loadedmetadata":           void,
	"loadend":                  void,
	"loadstart":                void,
	"lostpointercapture":       void,
	"messageerror":             void,
	"mousedown":                void,
	"mouseenter":               void,
	"mouseleave":               void,
	"mousemove":                void,
	"mouseout":                 void,
	"mouseover":                void,
	"mouseup":                  void,
	"noupdate":                 void,
	"nomatch":                  void,
	"notificationclick":        void,
	"orientationchange":        void,
	"pagehide":                 void,
	"pageshow":                 void,
	"pointercancel":            void,
	"pointerdown":              void,
	"pointerenter":             void,
	"pointerleave":             void,
	"pointerlockchange":        void,
	"pointerlockerror":         void,
	"pointermove":              void,
	"pointerout":               void,
	"pointerover":              void,
	"pointerup":                void,
	"popstate":                 void,
	"pushsubscriptionchange":   void,
	"ratechange":               void,
	"readystatechange":         void,
	"resourcetimingbufferfull": void,
	"selectstart":              void,
	"selectionchange":          void,
	"slotchange":               void,
	"soundend":                 void,
	"soundstart":               void,
	"speechend":                void,
	"speechstart":              void,
	"timeupdate":               void,
	"touchcancel":              void,
	"touchend":                 void,
	"touchenter":               void,
	"touchleave":               void,
	"touchmove":                void,
	"touchstart":               void,
	"transitionend":            void,
	"updateready":              void,
	"upgradeneeded":            void,
	"userproximity":            void,
	"versionchange":            void,
	"visibilitychange":         void,
	"voiceschanged":            void,
	"volumechange":             void,
	"vrdisplayconnected":       void,
	"vrdisplaydisconnected":    void,
	"vrdisplaypresentchange":   void,
	"click":                    void,
	"change":                   void,
	"input":                    void,
}
