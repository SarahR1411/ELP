module Main exposing (..)

import Browser
import Html exposing (Html, div, input, button, text)
import Html.Attributes exposing (value, type_, placeholder)
import Html.Events exposing (onInput, onClick)
import Svg exposing (Svg, svg)
import Svg.Attributes exposing (width, height, viewBox, style)
import Parser exposing (run)
import TcTurtleParser exposing (parseTurtleProgram)
import Turtle exposing (executeInstructions)
import View exposing (renderSvg)
import Html.Attributes as HtmlStyle exposing (style)

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
    div 
        [ HtmlStyle.style "display" "flex"
        , HtmlStyle.style "flex-direction" "column"
        , HtmlStyle.style "align-items" "center"
        , HtmlStyle.style "margin-top" "20px"
        ]
        [ div 
            [ HtmlStyle.style "margin-bottom" "10px" ]
            [ input 
                [ type_ "text"
                , onInput NewInput
                , value model.input
                , placeholder "Enter your TcTurtle command here..."
                , HtmlStyle.style "width" "400px"
                , HtmlStyle.style "padding" "10px"
                , HtmlStyle.style "border" "1px solid #ccc"
                , HtmlStyle.style "border-radius" "5px"
                ] 
                []
            , button 
                [ onClick Draw
                , HtmlStyle.style "margin-left" "10px"
                , HtmlStyle.style "padding" "10px"
                , HtmlStyle.style "border" "none"
                , HtmlStyle.style "background-color" "#007BFF"
                , HtmlStyle.style "color" "white"
                , HtmlStyle.style "border-radius" "5px"
                , HtmlStyle.style "cursor" "pointer"
                , HtmlStyle.style "font-weight" "bold"
                ] 
                [ text "Draw" ]
            ]
        , div 
            [ HtmlStyle.style "margin-top" "20px"
            , HtmlStyle.style "border" "2px solid #ddd"
            , HtmlStyle.style "width" "500px"
            , HtmlStyle.style "height" "500px"
            , HtmlStyle.style "display" "flex"
            , HtmlStyle.style "align-items" "center"
            , HtmlStyle.style "justify-content" "center"
            , HtmlStyle.style "background-color" "#f9f9f9"
            ]
            [ model.svg ]
        ]


-- MAIN ENTRY POINT
main =
    Browser.sandbox { init = init, update = update, view = view }
