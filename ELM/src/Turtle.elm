module Turtle exposing (Instruction(..), executeInstructions)

import Svg exposing (line)
import Svg.Attributes exposing (x1, y1, x2, y2, stroke, strokeWidth)

-- Turtle Instruction Type
type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)

type alias Turtle =
    { x : Float, y : Float, angle : Float }

moveTurtle : Instruction -> Turtle -> (Turtle, List (Svg.Svg msg))
moveTurtle instruction turtle =
    case instruction of
        Forward d ->
            let
                rad = degrees turtle.angle
                newX = turtle.x + toFloat d * cos rad
                newY = turtle.y - toFloat d * sin rad
                newTurtle = { turtle | x = newX, y = newY }
            in
            ( newTurtle, [ viewLine turtle newTurtle ] )

        Left d ->
            ( { turtle | angle = turtle.angle + toFloat d }, [] )

        Right d ->
            ( { turtle | angle = turtle.angle - toFloat d }, [] )

        Repeat n instructions ->
            let
                execute times (state, lines) =
                    if times == 0 then
                        ( state, lines )
                    else
                        let
                            (newState, newLines) = executeInstructions instructions state
                        in
                        execute (times - 1) (newState, lines ++ newLines)
            in
            execute n (turtle, [])

executeInstructions : List Instruction -> Turtle -> (Turtle, List (Svg.Svg msg))
executeInstructions instructions turtle =
    List.foldl
        (\instr (state, lines) ->
            let
                (newState, newLines) = moveTurtle instr state
            in
            ( newState, lines ++ newLines )
        )
        (turtle, [])
        instructions


viewLine : Turtle -> Turtle -> Svg.Svg msg
viewLine from to =
    line
        [ x1 (String.fromFloat from.x)
        , y1 (String.fromFloat from.y)
        , x2 (String.fromFloat to.x)
        , y2 (String.fromFloat to.y)
        , stroke "black"
        , strokeWidth "2"
        ]
        []
