module Main exposing (..)

import Browser
import Html exposing (Html, div, input, button, text)
import Html.Attributes exposing (value, type_, placeholder, style, min, max, step)
import Html.Events exposing (onInput, onClick)
import Svg exposing (Svg, svg)
import Svg.Attributes exposing (width, height, viewBox)
import Parser exposing (run)
import TcTurtleParser exposing (parseTurtleProgram)
import Turtle exposing (Instruction, Turtle, executeInstructions)
import View exposing (renderSvg)
import Time

-- MODEL
type alias Model =
    { input : String
    , svg : Html Msg
    , remainingInstructions : List Instruction
    , currentLines : List (Svg Msg)
    , currentTurtle : Turtle
    , isAnimating : Bool
    , animationSpeed : Int
    , penColor : String
    , penWidth : Int
    , penDown : Bool
    }

type Msg
    = NewInput String
    | Draw
    | StartAnimation
    | NextStep
    | AnimationTick
    | UpdateSpeed Int
    | UpdateColor String
    | UpdateWidth Int
    | TogglePen

init : Model
init =
    { input = ""
    , svg = renderSvg []
    , remainingInstructions = []
    , currentLines = []
    , currentTurtle = { x = 250, y = 250, angle = 0, penDown = True, penColor = "#007BFF", penWidth = 2 }
    , isAnimating = False
    , animationSpeed = 50
    , penColor = "#007BFF"
    , penWidth = 2
    , penDown = True
    }

-- UPDATE FUNCTION
update : Msg -> Model -> Model
update msg model =
    case msg of
        NewInput newText ->
            { model | input = newText }

        Draw ->
            case run parseTurtleProgram model.input of
                Ok parsedInstructions ->
                    let
                        initialTurtle = { x = 250, y = 250, angle = 0, penDown = model.penDown, penColor = model.penColor, penWidth = model.penWidth }
                        (finalTurtle, svgElements) = executeInstructions (expandRepeats parsedInstructions) initialTurtle
                    in
                    { model 
                        | svg = renderSvg svgElements
                        , currentTurtle = finalTurtle
                    }

                Err _ ->
                    { model | svg = text "Invalid syntax!" }

        StartAnimation ->
            case run parseTurtleProgram model.input of
                Ok parsedInstructions ->
                    { model
                        | remainingInstructions = expandRepeats parsedInstructions
                        , currentLines = []
                        , currentTurtle = { x = 250, y = 250, angle = 0, penDown = model.penDown, penColor = model.penColor, penWidth = model.penWidth }
                        , svg = renderSvg []
                        , isAnimating = True
                    }

                Err _ ->
                    { model | svg = text "Invalid syntax!" }

        NextStep ->
            if not model.isAnimating then
                model
            else
                case model.remainingInstructions of
                    [] ->
                        { model | isAnimating = False }

                    instr :: rest ->
                        let
                            (newTurtle, newLines) = executeInstructions [instr] model.currentTurtle
                        in
                        { model
                            | remainingInstructions = rest
                            , currentTurtle = newTurtle
                            , currentLines = model.currentLines ++ newLines
                            , svg = renderSvg (model.currentLines ++ newLines)
                        }

        AnimationTick ->
            if model.isAnimating then
                update NextStep model
            else
                model

        UpdateSpeed newSpeed ->
            { model | animationSpeed = newSpeed }

        UpdateColor color ->
            { model | penColor = color }

        UpdateWidth width ->
            { model | penWidth = width }

        TogglePen ->
            { model | penDown = not model.penDown }

-- HELPER FUNCTIONS
expandRepeats : List Instruction -> List Instruction
expandRepeats instructions =
    List.concatMap
        (\instr ->
            case instr of
                Turtle.Repeat n list ->
                    List.repeat n (expandRepeats list) |> List.concat

                other ->
                    [other]
        )
        instructions

exampleCommand : String -> String -> Html Msg
exampleCommand title cmd =
    div [ style "margin-bottom" "10px" ]
        [ div [ style "font-weight" "600", style "margin-bottom" "4px", style "color" "#2c3e50" ] [ text title ]
        , button
            [ onClick (NewInput cmd)
            , style "background-color" "#e8f4ff"
            , style "color" "#007BFF"
            , style "border" "1px solid #007BFF"
            , style "padding" "8px 12px"
            , style "border-radius" "6px"
            , style "cursor" "pointer"
            , style "width" "100%"
            , style "text-align" "left"
            ]
            [ text cmd ]
        ]

-- SUBSCRIPTIONS
subscriptions : Model -> Sub Msg
subscriptions model =
    if model.isAnimating then
        Time.every (toFloat model.animationSpeed) (always AnimationTick)
    else
        Sub.none

-- VIEW FUNCTION
view : Model -> Html Msg
view model =
    div 
        [ style "display" "flex"
        , style "justify-content" "center"
        , style "align-items" "flex-start"
        , style "min-height" "100vh"
        , style "padding-top" "40px"
        , style "background-color" "#f8f9fa"
        ]
        [ div 
            [ style "display" "flex"
            , style "gap" "40px"
            , style "max-width" "1200px"
            , style "width" "100%"
            , style "padding" "0 20px"
            ]
            [ div [ style "flex" "1" ]
                [ div 
                    [ style "border" "3px solid #e0e0e0"
                    , style "border-radius" "12px"
                    , style "overflow" "hidden"
                    , style "box-shadow" "0 4px 6px rgba(0, 0, 0, 0.1)"
                    , style "background-color" "white"
                    ]
                    [ model.svg ]
                , div 
                    [ style "margin-top" "20px"
                    , style "display" "flex"
                    , style "flex-direction" "column"
                    , style "gap" "15px"
                    ]
                    [ div 
                        [ style "display" "flex"
                        , style "gap" "10px"
                        , style "justify-content" "center"
                        ]
                        [ input
                            [ type_ "text"
                            , onInput NewInput
                            , value model.input
                            , placeholder "Enter Turtle commands (e.g., [Repeat 4 [Forward 50 Left 90]])"
                            , style "width" "100%"
                            , style "padding" "12px"
                            , style "border" "2px solid #007BFF"
                            , style "border-radius" "8px"
                            , style "font-size" "16px"
                            ]
                            []
                        ]
                    , div 
                        [ style "display" "flex"
                        , style "gap" "10px"
                        , style "justify-content" "center"
                        ]
                        [ button
                            [ onClick Draw
                            , style "padding" "12px 24px"
                            , style "background-color" "#007BFF"
                            , style "color" "white"
                            , style "border" "none"
                            , style "border-radius" "8px"
                            , style "cursor" "pointer"
                            , style "font-weight" "bold"
                            ]
                            [ text "Draw Instantly" ]
                        , button
                            [ onClick StartAnimation
                            , style "padding" "12px 24px"
                            , style "background-color" "#28A745"
                            , style "color" "white"
                            , style "border" "none"
                            , style "border-radius" "8px"
                            , style "cursor" "pointer"
                            , style "font-weight" "bold"
                            ]
                            [ text "Start Animation" ]
                        ]
                    , div 
                        [ style "display" "flex"
                        , style "flex-direction" "column"
                        , style "gap" "8px"
                        , style "background-color" "#f8f9fa"
                        , style "padding" "15px"
                        , style "border-radius" "8px"
                        ]
                        [ div 
                            [ style "color" "#6c757d"
                            , style "font-size" "14px"
                            , style "font-weight" "600"
                            ]
                            [ text "ANIMATION SPEED CONTROL" ]
                        , div 
                            [ style "display" "flex"
                            , style "align-items" "center"
                            , style "gap" "15px"
                            ]
                            [ input
                                [ type_ "range"
                                , Html.Attributes.min "10"
                                , Html.Attributes.max "250"
                                , step "10"
                                , value (String.fromInt model.animationSpeed)
                                , onInput (String.toInt >> Maybe.withDefault 50 >> UpdateSpeed)
                                , style "flex" "1"
                                , style "accent-color" "#007BFF"
                                ]
                                []
                            , div 
                                [ style "color" "#2c3e50"
                                , style "font-weight" "500"
                                , style "min-width" "80px"
                                ]
                                [ text (String.fromInt model.animationSpeed ++ "ms") ]
                            ]
                        ]
                    , div 
                        [ style "display" "flex"
                        , style "gap" "10px"
                        , style "align-items" "center"
                        ]
                        [ div 
                            [ style "display" "flex"
                            , style "flex-direction" "column"
                            , style "gap" "8px"
                            , style "flex" "1"
                            ]
                            [ div 
                                [ style "color" "#6c757d"
                                , style "font-size" "14px"
                                , style "font-weight" "600"
                                ]
                                [ text "PEN CONTROLS" ]
                            , div 
                                [ style "display" "flex"
                                , style "gap" "10px"
                                , style "align-items" "center"
                                ]
                                [ input
                                    [ type_ "color"
                                    , value model.penColor
                                    , onInput UpdateColor
                                    , style "height" "40px"
                                    , style "padding" "3px"
                                    , style "border" "1px solid #ddd"
                                    , style "border-radius" "6px"
                                    ]
                                    []
                                , input
                                    [ type_ "range"
                                    , Html.Attributes.min "1"
                                    , Html.Attributes.max "20"
                                    , value (String.fromInt model.penWidth)
                                    , onInput (String.toInt >> Maybe.withDefault 2 >> UpdateWidth)
                                    , style "flex" "1"
                                    , style "accent-color" model.penColor
                                    ]
                                    []
                                , button
                                    [ onClick TogglePen
                                    , style "padding" "10px 20px"
                                    , style "border" "none"
                                    , style "border-radius" "6px"
                                    , style "cursor" "pointer"
                                    , style "font-weight" "bold"
                                    , style "background-color" (if model.penDown then "#28a745" else "#dc3545")
                                    , style "color" "white"
                                    ]
                                    [ text (if model.penDown then "Pen is Down" else "Pen is Up") ]
                                ]
                            ]
                        ]
                    ]
                ]
            , div 
                [ style "width" "300px"
                , style "background-color" "white"
                , style "padding" "20px"
                , style "border-radius" "12px"
                , style "box-shadow" "0 4px 6px rgba(0, 0, 0, 0.05)"
                ]
                [ div 
                    [ style "color" "#2c3e50"
                    , style "font-size" "20px"
                    , style "font-weight" "600"
                    , style "margin-bottom" "10px"
                    ]
                    [ text "Example Commands" ]
                , div 
                    [ style "color" "#6c757d"
                    , style "font-size" "14px"
                    , style "margin-bottom" "20px"
                    ]
                    [ text "Click any example to use it in the input!" ]
                , exampleCommand "Square" "[Repeat 4 [Forward 50, Left 90]]"
                , exampleCommand "Triangle" "[Repeat 3 [Forward 60, Left 120]]"
                , exampleCommand "Hexagon" "[Repeat 6 [Forward 40, Left 60]]"
                , exampleCommand "Star" "[Repeat 5 [Forward 100, Right 144]]"
                , exampleCommand "Spiral" "[Repeat 36 [Forward 20, Left 10 ,Repeat 4 [Forward 10, Left 90]]]"
                , exampleCommand "Flower" "[Repeat 12 [Repeat 6 [Forward 30, Right 60], Right 30]]"
                , exampleCommand "Circle" "[Repeat 36 [Forward 10, Right 10]]"
                , exampleCommand "Sunburst" "[Repeat 24 [Forward 50, Right 170]]"
                , exampleCommand "Zigzag" "[Repeat 10 [Forward 30, Right 45, Forward 30, Left 90]]"
                ]
            ]
        ]


-- MAIN ENTRY POINT
main =
    Browser.element
        { init = \() -> (init, Cmd.none)
        , update = \msg model -> (update msg model, Cmd.none)
        , subscriptions = subscriptions
        , view = view
        }