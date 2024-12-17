# Widgets : Button

The button widget is the basic tappable interaction for an app.
A user tapping this will run the function passed to the constructor function. 

## Basic usage

Simply create a button using the `NewButton` constructor function, passing in a function
that should run when the button is tapped.

```
btn := widget.NewButton("Tap me", func() {})
```

If you want to use an icon in your button that is possible.
You can also set the label to "" if you want icon only!

```
btn := widget.NewButtonWithIcon("Home",
    theme.HomeIcon(), func() {})
```

## Disabled

A button can also be disabled so that it cannot be tapped:

```
btn := widget.NewButton("Tap me", func() {})
btn.Disable()
```

## Importance

You can change the colour / style of the button by setting its `Importance` value, like this:

```
btn := widget.NewButton("Danger!", func() {})
btn.Importance = widget.DangerImportance
```