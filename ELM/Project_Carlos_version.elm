module Project_Carlos_version exposing (..)


import Browser
import Html exposing (Html, div, button, text, canvas)
import Html.Attributes exposing (style, width, height)
import Html.Events exposing (onClick)
import Json.Encode as Encode


-- MODEL

type alias Model =
    { x : Float
    , y : Float
    , angle : Float
    , penDown : Bool
    , lines : List (Float, Float, Float, Float) -- (x1, y1, x2, y2)
    }

initialModel : Model
initialModel =
    { x = 300
    , y = 300
    , angle = 0
    , penDown = True
    , lines = []
    }


-- UPDATE

type Msg
    = Forward Float
    | Turn Float
    | PenUp
    | PenDown
    | Clear

update : Msg -> Model -> Model
update msg model =
    case msg of
        Forward distance ->
            let
                newX = model.x + distance * cos (degreesToRadians model.angle)
                newY = model.y - distance * sin (degreesToRadians model.angle)
            in
            if model.penDown then
                { model
                    | x = newX
                    , y = newY
                    , lines = (model.x, model.y, newX, newY) :: model.lines
                }
            else
                { model | x = newX, y = newY }

        Turn angleChange ->
            { model | angle = model.angle + angleChange }

        PenUp ->
            { model | penDown = False }

        PenDown ->
            { model | penDown = True }

        Clear ->
            { model | lines = [] }


-- VIEW

view : Model -> Html Msg
view model =
    div [ style "text-align" "center" ]
        [ div []
            [ button [ onClick (Forward 50) ] [ text "Forward" ]
            , button [ onClick (Turn 15) ] [ text "Turn Left" ]
            , button [ onClick (Turn -15) ] [ text "Turn Right" ]
            , button [ onClick PenUp ] [ text "Pen Up" ]
            , button [ onClick PenDown ] [ text "Pen Down" ]
            , button [ onClick Clear ] [ text "Clear" ]
            ]
        , canvasView model
        ]

canvasView : Model -> Html msg
canvasView model =
    Html.canvas
        [ width 600, height 600, style "border" "1px solid black" ]
        [ Encode.list (drawLines model.lines) ]


drawLines : List (Float, Float, Float, Float) -> List (Encode.Value)
drawLines lines =
    List.map drawLine lines


drawLine : (Float, Float, Float, Float) -> Encode.Value
drawLine ( x1, y1, x2, y2 ) =
    Encode.object
        [ ( "type", Encode.string "line" )
        , ( "x1", Encode.float x1 )
        , ( "y1", Encode.float y1 )
        , ( "x2", Encode.float x2 )
        , ( "y2", Encode.float y2 )
        , ( "color", Encode.string "black" )
        , ( "width", Encode.float 2 )
        ]


-- HELPER FUNCTIONS

degreesToRadians : Float -> Float
degreesToRadians degrees =
    degrees * pi / 180


-- MAIN

main : Program () Model Msg
main =
    Browser.sandbox { init = initialModel, update = update, view = view }
