module View exposing (renderSvg)

import Svg exposing (Svg, svg)
import Svg.Attributes exposing (..)
import Html exposing (Html)


renderSvg : List (Svg msg) -> Html msg
renderSvg lines =
    svg [ width "500", height "500", viewBox "0 0 500 500" ] lines
