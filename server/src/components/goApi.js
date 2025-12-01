import { useEffect, useState } from "react"



const usePreConfig = () => {
    const [data, setData] = useState([])
    useEffect(() => {
        fetch("/api/get/preConfig").then(resp => resp.json()).then(data => setData(data))
        // setInterval(() => {//其实没必要，运行期间不会变
        //     fetch("/api/get/macros").then(resp => resp.json()).then(data => setData(data))
        // }, 3000)
        return () => { }
    }, [])
    return data
}


const useMacros = () => {
    const [data, setData] = useState([])
    useEffect(() => {
        fetch("/api/get/macros").then(resp => resp.json()).then(data => setData(data))
        // setInterval(() => {//其实没必要，运行期间不会变
        //     fetch("/api/get/macros").then(resp => resp.json()).then(data => setData(data))
        // }, 3000)
        return () => { }
    }, [])
    return data
}


const useMouseConfig = () => {
    const [data, setData] = useState([])
    useEffect(() => {
        fetch("/api/get/mouse").then(resp => resp.json()).then(data => setData(data))
        // setInterval(() => {
        //     fetch("/api/get/mouse").then(resp => resp.json()).then(data => setData(data))
        // }, 1000)
        return () => { }
    }, [])
    const setMouse = async (key,value) => {
        await fetch(`/api/set/mouse?key=${key}&value=${value}` )
        fetch("/api/get/mouse").then(resp => resp.json()).then(data => setData(data))
    }   
    return [data,setMouse]
}

const useKeyboardConfig = () => {
    const [data, setData] = useState([])
    useEffect(() => {
        fetch("/api/get/keyboard").then(resp => resp.json()).then(data => setData(data))
        // setInterval(() => {
        //     fetch("/api/get/keyboard").then(resp => resp.json()).then(data => setData(data))
        // }, 1000)
        return () => { }
    }, [])
    const setKeyboard = async (key,value) => {
        await fetch(`/api/set/keyboard?key=${key}&value=${value}` )
        fetch("/api/get/keyboard").then(resp => resp.json()).then(data => setData(data))
    } 
    return [data,setKeyboard]
}

export { useMacros, useMouseConfig, useKeyboardConfig,usePreConfig }