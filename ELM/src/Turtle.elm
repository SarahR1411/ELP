module Turtle exposing (Instruction(..), executeInstructions, Turtle)

-- Import the line function from the Svg module to draw lines
import Svg exposing (line)
-- Import specific attributes from Svg.Attributes to control line properties like coordinates and stroke
import Svg.Attributes exposing (x1, y1, x2, y2, stroke, strokeWidth)

-- Define the Instruction type, which represents different turtle movements
type Instruction
    = Forward Int   -- Move the turtle forward by a specified number of units (Int)
    | Left Int    -- Turn the turtle left by a specified number of degrees (Int)
    | Right Int  -- Turn the turtle right by a specified number of degrees (Int)
    | Repeat Int (List Instruction) -- Repeat a sequence of instructions a specified number of times

-- Define the Turtle type as a record with properties like position, angle, pen state, etc.
type alias Turtle =
    { x : Float -- X position of the turtle
    , y : Float -- Y position of the turtle
    , angle : Float -- Angle of the turtle in degrees
    , penDown : Bool    -- Whether the pen is down (drawing) or up (not drawing)
    , penColor : String -- Color of the pen
    , penWidth : Int        -- Width of the pen stroke
    }

-- The moveTurtle function moves the turtle based on a single Instruction
moveTurtle : Instruction -> Turtle -> (Turtle, List (Svg.Svg msg))
moveTurtle instruction turtle =
    case instruction of
        -- Move forward by a certain distance
        Forward d ->
            let
                -- Convert the angle to radians for trigonometric calculations
                rad = degrees turtle.angle

                -- Calculate the new X and Y positions based on the distance and angle
                newX = turtle.x + toFloat d * cos rad
                newY = turtle.y - toFloat d * sin rad

                -- Create a new Turtle record with updated position
                newTurtle = { turtle | x = newX, y = newY }
            in
            -- If the pen is down, generate a line from the previous position to the new one
            if turtle.penDown then
                ( newTurtle, [ viewLine turtle newTurtle turtle.penColor turtle.penWidth ] )
            else
                ( newTurtle, [] )

        -- Turn the turtle left by a certain angle
        Left d ->
            ( { turtle | angle = turtle.angle + toFloat d }, [] )
        -- Turn the turtle right by a certain angle
        Right d ->
            ( { turtle | angle = turtle.angle - toFloat d }, [] )

        -- Repeat a set of instructions (this is a placeholder for future handling)
        Repeat _ _ ->
            ( turtle, [] )

-- The executeInstructions function takes a list of Instructions and a Turtle, and processes them sequentially
executeInstructions : List Instruction -> Turtle -> (Turtle, List (Svg.Svg msg))
executeInstructions instructions turtle =
    List.foldl
    -- Fold function that processes each instruction and accumulates the turtle's state and lines drawn
        (\instr (state, lines) ->
            let
            -- Process the current instruction and get the updated state and lines
                (newState, newLines) = moveTurtle instr state
            in
            -- Accumulate the new state and the lines drawn so far
            ( newState, lines ++ newLines )
        )
        (turtle, [])    -- Initial state is the current turtle and an empty list of lines
        instructions

-- The viewLine function generates an SVG line element that represents the movement from one turtle position to another
viewLine : Turtle -> Turtle -> String -> Int -> Svg.Svg msg
viewLine from to color width =
    line
        [ -- Set the starting (x1, y1) position for the line from the "from" turtle
        x1 (String.fromFloat from.x)
        , y1 (String.fromFloat from.y)

        -- Set the ending (x2, y2) position for the line at the "to" turtle's new position
        , x2 (String.fromFloat to.x)
        , y2 (String.fromFloat to.y)

        -- Set the color of the line stroke
        , stroke color

        -- Set the width of the line
        , strokeWidth (String.fromInt width)
        ]
        []  -- No children elements, just the line itself