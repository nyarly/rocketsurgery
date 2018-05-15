# Rocket Surgery

This is a library package,
intended for building tools that transform Go source code.
It's modelled in part on the general pattern of
a templating library,
in that you build a code template
for the overall pattern of the transformation target,
parse input code into and AST,
and then use the parsed code as a context as input
to produce the target AST,
which is then rendered into files.

Rather than operate at a textual level,
`rocketsurgery` works with the Go AST.
The benefits of this are that you can manipulate units of code,
and be sure that your output will be syntactic.
It also ensures, for instance, than names won't collide by accident.
The vision overall has been of a kind of pragmatic hygenic macro system.

One challenge in using the Go AST parsing engines
is in adding any kind of templating logic
to the template sources themselves -
there's just not much semantic room to make those ergonomic.
Instead, the principle is that this is a toolkit,
and the logic to transform the templates will accompany them.

n.B. this package is absolutely in an early experimental stage.
Good luck with it!
