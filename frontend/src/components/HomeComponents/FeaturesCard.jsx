import React from 'react'
import "../CSS/FeaturesCard.css"
// import Notetaking from "../Assets/notetaking.svg"


function Features_card(props) {
    return (
        <div className='card'>
            <img src={props.img} alt="" />
            <h3>{props.title}</h3>
        </div>
    )
}

export default Features_card