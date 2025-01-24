-- This is just an example to understand the structure of an Elm program!!!

module main exposing( main )

import Browser exposing(Document, UrlRequest, UrlResponse)
import Html exposing( Html, div, input, text )
import Browser.Navigation exposing (Key)
import Url exposing (Url)
-- Imports modules and exposes specific functions/types:
    -- Browser: For building applications, working with the browser, and routing.
    -- Html: For constructing the DOM and rendering HTML.
    -- Browser.Navigation: Provides functionality for navigation and handling URLs.
    -- Url: For working with URLs.



type alias Flags = 
    {}
    -- Defines an alias Flags as an empty record. Flags are passed to the application from JavaScript.


type alias Model = 
    { name : String --just a placeholder
    }
    -- Defines a Model with a single field name of type String. This will represent the application's state.


type Msg
    = NoOp --justa a placeholder
    -- Defines a message type Msg with a single constructor NoOp. Messages are used to signal events or changes in the application.


init : Flags -> Url -> Key -> (Model, Cmd Msg) --First function called when prog ran, gets called with any flags, url, and key
init flags url key = 
    let
        _ =
            Debug.log "url" url
    in
    (Model "Jack", Cmd.none)
     -- Logs the URL using Debug.log.
    -- Initializes the model with the name "Jack" and specifies no commands (Cmd.none).


view : Model -> Document Msg
-- Defines how the application’s Model is transformed into a view.
-- Returns a Document Msg, which includes:
-- title: The document's title.
-- body: A list of HTML elements
view model = 
    {
        tittle = "Distinctly Average"
        , body = 
            [Html.p[] [Html.text "Hello, World!"]
            ]
    }
    -- Specifies the title "Distinctly Average".
    -- Renders a single paragraph with the text "Hello, World!".


onUrlRequest : UrlRequest -> Msg
onUrlRequest urlRequest = 
    NoOp
-- Handles URL requests, returning NoOp for now.


onUrlchange : Url -> Msg
onUrlChange url = 
    NoOp   
-- Handles changes to the URL, also returning NoOp.


update : Msg -> Model -> (Model, Cmd Msg)
update msg model = 
    case msg of
        NoOp -> 
            (model, Cmd.none)
-- Handles messages (Msg) to update the Model.
-- In this case, NoOp leaves the model unchanged and specifies no side effects.


main = 
    Browser.application
    {
        init = init
        , view = view
        , update = update
        , subscription = \_ -> Sub.none
        , onUrlRequest = onUrlRequest
        , onUrlChange = onUrlChange
    }
-- Defines the main application using Browser.application.
-- Specifies:
-- init: Initialization function.
-- view: Function to generate the view.
-- update: Function to update the model.
-- subscriptions: No subscriptions (Sub.none).
-- onUrlRequest and onUrlChange: Handle URL requests and changes.


