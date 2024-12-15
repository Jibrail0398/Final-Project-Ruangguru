function ChatEnd(props){
    return(
        <div className="chat chat-end">
            <div className="chat-bubble chat-bubble-warning ml-4 mt-8">{props.message}</div>
        </div>
    )
}

export default ChatEnd;