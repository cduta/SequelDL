package backend

type Options struct {
  WindowTitle     string
  WindowWidth     int64
  WindowHeight    int64
  DefaultFont     string 
  DefaultFontSize int64 
  FPS             int64
  ShowFPS         bool
}

func (handle *Handle) QueryOptions() (Options, error) {
	var (
		err      error 
		row     *Row
		options  Options
	)

	row, err = handle.QueryRow(`
SELECT window_title.value,
       window_width.value,
       window_height.value,
       default_font.value,
       default_font_size.value,
       fps.value,
       show_fps.value
FROM integer_options AS window_width,
     integer_options AS window_height,
     integer_options AS default_font_size,
     integer_options AS fps,
     text_options    AS window_title,
     text_options    AS default_font,
     boolean_options AS show_fps
WHERE window_title.name      = 'window title'
AND   window_width.name      = 'window width'
AND   window_height.name     = 'window height'
AND   default_font.name      = 'default font'
AND   default_font_size.name = 'default font size'
AND   fps.name               = 'fps'
AND   show_fps.name          = 'show fps';
	`);
  if err != nil {
    return options, err
  }

  if row != nil {
    err = row.Scan(
    	&options.WindowTitle,
      &options.WindowWidth,
      &options.WindowHeight,
      &options.DefaultFont,
      &options.DefaultFontSize,
      &options.FPS,
      &options.ShowFPS)
  } 

	return options, err
}