module View exposing (renderSvg)

-- Import the Svg module, which provides the Svg type and the svg function
import Svg exposing (Svg, svg)

-- Import the Svg.Attributes module to set attributes for the SVG elements
import Svg.Attributes exposing (..)

-- Import the Html module to return HTML as a result
import Html exposing (Html)

-- Define the renderSvg function, which takes a list of SVG elements and returns HTML
renderSvg : List (Svg msg) -> Html msg
renderSvg lines =
    -- The svg function generates the <svg> element.
    -- The first argument is a list of attributes for the <svg> element.
    svg [ width "500", height "500", viewBox "0 0 500 500" ] lines
    -- Set the width of the SVG to 500 units
    -- Set the height of the SVG to 500 units
    -- Set the viewBox attribute to "0 0 500 500"
     -- The list of SVG elements (e.g., circles, paths) to include inside the <svg>
