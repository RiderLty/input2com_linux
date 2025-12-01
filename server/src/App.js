
import { CssBaseline, Grid, useEventCallback } from '@mui/material';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import './App.css';

import KeyContainer from './components/keyContainer';
import MouseContainer from "./components/mouseContainer"
import MacroPanel from './components/macroPanel';
import PreConfigPanel from './components/preConfigPanel';
import { DndProvider } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import LongClickMenu from './components/LongClickMenu';
import CenterFocusWeakIcon from '@material-ui/icons/CenterFocusWeak';
import { useKeyboardConfig, useMacros, useMouseConfig,usePreConfig } from './components/goApi';
import ClearAllIcon from '@material-ui/icons/ClearAll';



const Keyboard = ({ setKeyboardConfig, onKeyboardClick }) => {
  return <div className="keyboard full-size">
    {/* from https://github.com/Mostafa-Abbasi/KeyboardTester */}
    <section className="function region">
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="escape" className="key escape">ESC</KeyContainer>
      <div className="empty-space-between-keys" aria-hidden="true"></div>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f1" className="key f1">F1</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f2" className="key f2">F2</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f3" className="key f3">F3</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f4" className="key f4">F4</KeyContainer>
      <div className="empty-space-between-keys" aria-hidden="true"></div>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f5" className="key f5">F5</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f6" className="key f6">F6</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f7" className="key f7">F7</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f8" className="key f8">F8</KeyContainer>
      <div className="empty-space-between-keys" aria-hidden="true"></div>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f9" className="key f9">F9</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f10" className="key f10">F10</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f11" className="key f11">F11</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="f12" className="key f12">F12</KeyContainer>
    </section>

    <section className="system-control region">
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="printscreen" className="key printscreen key--accent-color">Prt Sc</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="scrolllock" className="key scrolllock key--accent-color">Scr Lk</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="pause" className="key pause key--accent-color">Pause</KeyContainer>
    </section>

    <section className="typewriter region">
      <div className="first-row">
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="backquote" className="key backquote key--sublegend key--accent-color">
          <span>~</span> <span>`</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit1" className="key digit1">1</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit2" className="key digit2">2</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit3" className="key digit3">3</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit4" className="key digit4">4</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit5" className="key digit5">5</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit6" className="key digit6">6</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit7" className="key digit7">7</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit8" className="key digit8">8</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit9" className="key digit9">9</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="digit0" className="key digit0">0</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="minus" className="key minus key--sublegend">
          <span>{"−"}</span> <span>{"‐"}</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="equal" className="key equal key--sublegend">
          <span>{"+"}</span><span>{"="}</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="backspace" className="key backspace key--accent-color">Backspace</KeyContainer>
      </div>

      <div className="second-row">
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="tab" className="key tab key--accent-color">Tab</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyq" className="key keyq">Q</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyw" className="key keyw">W</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keye" className="key keye">E</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyr" className="key keyr">R</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyt" className="key keyt">T</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyy" className="key keyy">Y</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyu" className="key keyu">U</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyi" className="key keyi">I</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyo" className="key keyo">O</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyp" className="key keyp">P</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="bracketleft" className="key bracketleft key--sublegend">
          <span>{"{"}</span> <span>{"["}</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="bracketright" className="key bracketright key--sublegend">
          <span>{"}"}</span> <span>{"]"}</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="backslash" className="key backslash key--sublegend key--accent-color">
          <span>{"|"}</span><span>{"\\"}</span>
        </KeyContainer>
      </div>

      <div className="third-row">
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="capslock" className="key capslock key--accent-color">Caps</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keya" className="key keya">A</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keys" className="key keys">S</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyd" className="key keyd">D</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyf" className="key keyf">F</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyg" className="key keyg">G</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyh" className="key keyh">H</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyj" className="key keyj">J</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyk" className="key keyk">K</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyl" className="key keyl">L</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="semicolon" className="key semicolon key--sublegend">
          <span>{":"}</span> <span>{";"}</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="quote" className="key quote key--sublegend">
          <span>&quot;</span> <span>&apos;</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="enter" className="key enter key--accent-color">
          <span>{"⟵"}</span>
        </KeyContainer>
      </div>

      <div className="fourth-row">
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="shiftleft" className="key shiftleft key--accent-color">Shift</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyz" className="key keyz">Z</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyx" className="key keyx">X</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyc" className="key keyc">C</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyv" className="key keyv">V</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyb" className="key keyb">B</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keyn" className="key keyn">N</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="keym" className="key keym">M</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="comma" className="key comma key--sublegend">
          <span>{"<"}</span> <span>{","}</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="period" className="key period key--sublegend">
          <span>{">"}</span> <span>{"."}</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="slash" className="key slash key--sublegend">
          <span>{"?"}</span> <span>{"/"}</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="shiftright" className="key shiftright key--accent-color">Shift</KeyContainer>
      </div>

      <div className="fifth-row">
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="controlleft" className="key controlleft key--accent-color">Ctrl</KeyContainer>
        <KeyContainer
          code="metaleft" className="key metaleft osleft key--accent-color"
          aria-label="Left Windows key"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 4875 4875"
            aria-label="Left Windows key icon"
          >
            <path
              d="M0 0h2311v2310H0zm2564 0h2311v2310H2564zM0 2564h2311v2311H0zm2564 0h2311v2311H2564"
            />
          </svg>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="altleft" className="key altleft key--accent-color">Alt</KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="space" className="key space key--accent-color" aria-label="Space">
          <span>{'____'}</span>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="altright" className="key altright key--accent-color">Alt</KeyContainer>
        <KeyContainer
          code="metaright" className="key metaright osright key--accent-color"
          aria-label="Right Windows key"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 4875 4875"
            aria-label="Right Windows key icon"
          >
            <path
              d="M0 0h2311v2310H0zm2564 0h2311v2310H2564zM0 2564h2311v2311H0zm2564 0h2311v2311H2564"
            />
          </svg>
        </KeyContainer>
        <KeyContainer
          code="contextmenu" className="key contextmenu key--accent-color"
          aria-label="Context Menu key"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 100 100"
            aria-label="Context Menu key icon"
          >
            <rect
              x="10"
              y="10"
              width="80"
              height="80"
              rx="2"
              ry="2"
              strokeWidth="7"
            />
            <rect
              x="25"
              y="30"
              width="50"
              height="4.5"
              rx="2"
              ry="2"
              strokeWidth="4.5"
            />
            <rect
              x="25"
              y="47.5"
              width="50"
              height="4.5"
              rx="2"
              ry="2"
              strokeWidth="4.5"
            />
            <rect
              x="25"
              y="65"
              width="50"
              height="4.5"
              rx="2"
              ry="2"
              strokeWidth="4.5"
            />
          </svg>
        </KeyContainer>
        <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="controlright" className="key controlright key--accent-color">Ctrl</KeyContainer>
      </div>
    </section>

    <section className="navigation region">
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="insert" className="key insert key--accent-color">Insert</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="home" className="key home key--accent-color">Home</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="pageup" className="key pageup key--accent-color">Pg Up</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="delete" className="key delete key--accent-color">Delete</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="end" className="key end key--accent-color">End</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="pagedown" className="key pagedown key--accent-color">Pg Dn</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="arrowup" className="key arrowup" aria-label="Up Arrow key">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          enableBackground="new 0 0 32 32"
          viewBox="0 0 32 32"
          aria-label="Up Arrow"
        >
          <path
            d="M18.221,7.206l9.585,9.585c0.879,0.879,0.879,2.317,0,3.195l-0.8,0.801c-0.877,0.878-2.316,0.878-3.194,0  l-7.315-7.315l-7.315,7.315c-0.878,0.878-2.317,0.878-3.194,0l-0.8-0.801c-0.879-0.878-0.879-2.316,0-3.195l9.587-9.585  c0.471-0.472,1.103-0.682,1.723-0.647C17.115,6.524,17.748,6.734,18.221,7.206z"
          />
        </svg>
      </KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="arrowleft" className="key arrowleft" aria-label="Left Arrow key">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          enableBackground="new 0 0 32 32"
          viewBox="0 0 32 32"
          aria-label="Left Arrow"
        >
          <path
            d="M7.701,14.276l9.586-9.585c0.879-0.878,2.317-0.878,3.195,0l0.801,0.8c0.878,0.877,0.878,2.316,0,3.194  L13.968,16l7.315,7.315c0.878,0.878,0.878,2.317,0,3.194l-0.801,0.8c-0.878,0.879-2.316,0.879-3.195,0l-9.586-9.587  C7.229,17.252,7.02,16.62,7.054,16C7.02,15.38,7.229,14.748,7.701,14.276z"
          />
        </svg>
      </KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="arrowdown" className="key arrowdown" aria-label="Down Arrow key">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          enableBackground="new 0 0 32 32"
          viewBox="0 0 32 32"
          aria-label="Down Arrow"
        >
          <path
            d="M14.77,23.795L5.185,14.21c-0.879-0.879-0.879-2.317,0-3.195l0.8-0.801c0.877-0.878,2.316-0.878,3.194,0  l7.315,7.315l7.316-7.315c0.878-0.878,2.317-0.878,3.194,0l0.8,0.801c0.879,0.878,0.879,2.316,0,3.195l-9.587,9.585  c-0.471,0.472-1.104,0.682-1.723,0.647C15.875,24.477,15.243,24.267,14.77,23.795z"
          />
        </svg>
      </KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="arrowright" className="key arrowright" aria-label="Right Arrow key">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          enableBackground="new 0 0 32 32"
          viewBox="0 0 32 32"
          aria-label="Right Arrow"
        >
          <path
            d="M24.291,14.276L14.705,4.69c-0.878-0.878-2.317-0.878-3.195,0l-0.8,0.8c-0.878,0.877-0.878,2.316,0,3.194  L18.024,16l-7.315,7.315c-0.878,0.878-0.878,2.317,0,3.194l0.8,0.8c0.878,0.879,2.317,0.879,3.195,0l9.586-9.587  c0.472-0.471,0.682-1.103,0.647-1.723C24.973,15.38,24.763,14.748,24.291,14.276z"
          />
        </svg>
      </KeyContainer>
    </section>

    <section className="numpad region">
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numlock" className="key numlock">NumLk</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpadKeyContaineride" className="key numpadKeyContaineride">/</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpadmultiply" className="key numpadmultiply">&times;</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpadsubtract" className="key numpadsubtract">&minus;</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad7" className="key numpad7">7</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad8" className="key numpad8">8</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad9" className="key numpad9">9</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpadadd" className="key numpadadd">{"+"}</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad4" className="key numpad4">4</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad5" className="key numpad5">5</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad6" className="key numpad6">6</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad1" className="key numpad1">1</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad2" className="key numpad2">2</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad3" className="key numpad3">3</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpadenter" className="key numpadenter">Enter</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpad0" className="key numpad0">0</KeyContainer>
      <KeyContainer setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} code="numpaddecimal" className="key numpaddecimal">&middot;</KeyContainer>
    </section>
  </div>
}




function App() {

  const [dark, setDark] = useState(false);

  useEffect(() => {
    if (typeof window === 'undefined' || !window.matchMedia) return;
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    setDark(mediaQuery.matches);
    document.querySelector('meta[name="theme-color"]').setAttribute('content', mediaQuery.matches ? '#303030' : '#ECEFF1')
    const handler = (e) => setDark(e.matches);
    mediaQuery.addEventListener('change', handler);
    return () => mediaQuery.removeEventListener('change', handler);
  }, []);

  const theme = createTheme({
    components: {
      MuiCssBaseline: {
        styleOverrides: {
          body: {
            backgroundColor: dark ? '#303030' : "#ECEFF1",
            "a: link": {
              color: dark ? '#00796b' : "#00796b",
            },
            "a: visited": {
              color: dark ? '#00796b' : "#00796b",
            },
            "a: active": {
              color: dark ? '#00796b' : "#00796b",
            },
          },
        },
      },
    },

    palette: dark ? {

      primary: {
        main: "#00796B",
      },
      secondary: {
        main: "#d90051",
      },
      background: {
        main: "#303030",
        secondary: "#fefefe",
        mainCard: "#fdfdfd",
        read: "#9E9E9E",
        readHover: "#BDBDBD",
        tag: "#00796b",
        tagHover: "#009688",
        pageShadow: "8px 8px 16px #c4c4c4,-8px -8px 16px #ffffff"
      },
      iconButton: {
        main: "#000000",
        disabled: "#9e9e9e",
      },
      button: {
        tag: {
          type: {
            main: "#4a4a4a",
            hover: "#646464"
          },
          value: {
            main: "#4a4a4a",
            hover: "#646464"
          },
          text: "#ffffff",
        },
        iconFunction: {
          main: "#ffffff",
          disabled: "#9e9e9e",
          text: "#ffffff",
          process: "#d90051",
        },
        macro: {
          main: "#4a4a4a",
          hover: "#646464",
          process: "#d90051",
          buffer: "e0e0e0",
          text: "#ffffff"
        },
        mouse: {
          main: "#4a4a4a",
          hover: "#646464",
          process: "#d90051",
          buffer: "e0e0e0",
          text: "#ffffff"
        },
        loadMore: {
          main: "#303030",
          hover: "#646464",
          text: "#ffffff"
        },
        galleryCard: {
          main: "#212121",
        },
        commentCard: {
          main: "#303030",
          hover: "#303030",
          text: "#ffffff"
        }
      },
      text: {
        primary: "#ffffff",
        secondary: "#dddddd",
        disabled: "#9e9e9e",
      },
      page: {
        background: "#303030",
        shadow: "8px 8px 16px #252525,-8px -8px 16px #3b3b3b"
      },
      search: {
        color: "#3a3a3a",
        text: "#ffffff",
        split: "#757575"
      }
    }
      :
      {
        primary: {
          main: "#d90051",
        },
        secondary: {
          main: "#00796B",
        },
        background: {
          main: "#ECEFF1",
          secondary: "#fefefe",
          mainCard: "#fdfdfd",
          read: "#9E9E9E",
          readHover: "#BDBDBD",
          tag: "#00796b",
          tagHover: "#009688",
          pageShadow: "8px 8px 16px #c4c4c4,-8px -8px 16px #ffffff"
        },
        iconButton: {
          main: "#000000",
          disabled: "#9e9e9e",
        },

        button: {
          tag: {
            type: {
              main: "#C2185B",
              hover: "#E91E63"
            },
            value: {
              main: "#00796b",
              hover: "#009688"
            },
            text: "#ffffff",
          },
          iconFunction: {
            main: "#000000",
            disabled: "#9e9e9e",
            text: "#ffffff",
            process: "#d90051",
          },
          macro: {
            main: "#9e9e9e",
            hover: "#bdbdbd",
            process: "#d90051",
            buffer: "e0e0e0",
            text: "#000000"
          },
          mouse: {
            main: "#9e9e9e",
            hover: "#bdbdbd",
            process: "#d90051",
            buffer: "e0e0e0",
            text: "#000000"
          },
          loadMore: {
            main: "#ECEFF1",
            hover: "#eeeeee",
            text: "#000000"
          },
          galleryCard: {
            main: "#ffffff",
          },
          commentCard: {
            main: "#ECEFF1",
            hover: "#ECEFF1",
            text: "#000000"
          }
        },
        text: {
          primary: "#000000",
          secondary: "#757575",
          disabled: "#9e9e9e",
        },
        page: {
          background: "#ECEFF1",
          shadow: "8px 8px 16px #c4c4c4,-8px -8px 16px #ECEFF1"
        },
        search: {
          color: "#eeeeee",
          text: "#000000",
          split: "#3a3a3a"
        }
      }
  });

  const nowSelectkey = useRef(["mouse", 1])
  const nowSelectValue = useRef(false)


  const [pos, setPos] = useState([-1, -1])//点击弹出
  const [longClickItems, setLongClickItems] = useState([
  ])
  const [longClickedName, setLongClickedName] = useState("")//弹出的标题
  const macros = useMacros()
  const preConfig = usePreConfig()
  const [mouseConfig, setMouseConfig] = useMouseConfig()
  const [keyboardConfig, setKeyboardConfig] = useKeyboardConfig()

  const onMouseClick = (code, x, y) => {
    nowSelectkey.current = ["mouse", code]
    setPos([x, y])
  }
  const onKeyboardClick = (code, x, y) => {
    nowSelectkey.current = ["keyboard", code]
    setPos([x, y])
  }

  useEffect(() => {
    setLongClickItems(old => {
      return [{
        text: "清除",
        onClick: () => {
          if (nowSelectkey.current[0] === "mouse") {
            setMouseConfig(nowSelectkey.current[1], "CLEAR_FUNCTION")
          } else if (nowSelectkey.current[0] === "keyboard") {
            setKeyboardConfig(nowSelectkey.current[1], "CLEAR_FUNCTION")
          }
        },
        icon: <ClearAllIcon />
      }, ...Object.keys(macros).map(key => {
        return {
          text: macros[key]["name"] + " : " + macros[key]["description"],
          onClick: () => {
            if (nowSelectkey.current[0] === "mouse") {
              setMouseConfig(nowSelectkey.current[1], key)
            } else if (nowSelectkey.current[0] === "keyboard") {
              setKeyboardConfig(nowSelectkey.current[1], key)
            }
          },
          icon: <CenterFocusWeakIcon />
        }
      })]
    })
    console.log(macros)
  }, [macros])

  return (
    <DndProvider backend={HTML5Backend}>
      <ThemeProvider theme={theme}>
        <LongClickMenu
          pos={pos}
          setPos={setPos}
          items={longClickItems}
          title={longClickedName}
        />
        <CssBaseline />
        <div id='mainContainer' style={{ backgroundColor: theme.palette.page.background, width: "100%" }}    >
          <Grid
            container
            direction="column"
            justifyContent="center"
            alignItems="center"
          >
            <Grid
              item
              container
              direction="row"
              justifyContent="center"
              alignItems="center"
              xs={12}
            >
              <MouseContainer setMouseConfig={setMouseConfig} onMouseClick={onMouseClick} />
              <MacroPanel macros={macros} />
            </Grid>
            <Grid
              item
              container
              direction="row"
              justifyContent="center"
              alignItems="center"
              xs={12}
            >
              <PreConfigPanel preConfig={preConfig}  />
            </Grid>
            <Grid item xs={12} style={{
              width: "100%"
            }}>
              <Keyboard setKeyboardConfig={setKeyboardConfig} onKeyboardClick={onKeyboardClick} />
            </Grid>
          </Grid>
        </div >
      </ThemeProvider>
    </DndProvider>
  );
}




export default App;
