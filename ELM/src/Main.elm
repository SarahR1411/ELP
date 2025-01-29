module Main exposing (..)

import Browser
import Html exposing (Html, div, input, button, text)
import Html.Attributes exposing (value, type_)
import Html.Events exposing (onInput, onClick)
import Svg exposing (svg)
import Parser exposing (run)
import TcTurtleParser exposing (parseTurtleProgram)
import Turtle exposing (executeInstructions)
import View exposing (renderSvg)

-- MODEL
type alias Model =
    { input : String
    , svg : Html Msg
    }

type Msg
    = NewInput String
    | Draw

init : Model
init =
    { input = ""
    , svg = text ""
    }

-- UPDATE FUNCTION
update : Msg -> Model -> Model
update msg model =
    case msg of
        NewInput newText ->
            { model | input = newText }

        Draw ->
            let
                parsedResult = run parseTurtleProgram model.input
            in
            case parsedResult of
                Ok parsedInstructions ->
                    let
                        (_, svgElements) = executeInstructions parsedInstructions { x = 250, y = 250, angle = 0 }
                    in
                    { model | svg = renderSvg svgElements }

                Err err ->
                    { model | svg = text ("Parsing error: " ++ Debug.toString err) }


-- VIEW FUNCTION
view : Model -> Html Msg
view model =
    div []
        [ input [ type_ "text", onInput NewInput, value model.input ] []
        , button [ onClick Draw ] [ text "Draw" ]
        , div [] [ model.svg ]
        ]

-- MAIN ENTRY POINT
main =
    Browser.sandbox { init = init, update = update, view = view }
