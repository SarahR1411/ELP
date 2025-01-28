module Buttons exposing (..)
import Browser
import Html exposing (Html, div, button, text)
import Html.Events exposing (onClick)

-- Main
main = Browser.sandbox { init = init, update = update, view = view }

-- Model
type alias Model = { count : Int }
init : Model
init = { count = 0 }

-- Update
type Msg
    = Increment
    | Decrement
update : Msg -> Model -> Model
update msg model =
    case msg of
        Increment ->
            { model | count = model.count + 1 }
        Decrement ->
            { model | count = model.count - 1 }

-- View
view : Model -> Html Msg
view model = 
    div []
        [ button [ onClick Increment ] [ text "+" ]
        , div [] [ text (String.fromInt model.count) ]
        , button [ onClick Decrement ] [ text "-" ]
        ]