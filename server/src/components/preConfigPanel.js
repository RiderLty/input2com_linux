import { Button, Grid, Paper, Typography } from "@mui/material";
import { DndSender, ItemTypes } from "./dnd";

const PreConfigPanel = ({ preConfig }) => {
    return <div style={{
        // height: "261px",
        position: "relative",
        padding: "24px",
        display: "inline-flex",
    }}>
        <Grid
            container
            direction="row"
            justifyContent="center"
            alignItems="flex-start"
            spacing={"5px"}
            sx={{
                width: "80vw",
                height: "100%",
            }}
        >
            {Object.keys(preConfig).map((value) => (
                <Grid key={value} item xs={2} >
                    <Button
                        sx={{
                            width: "100%",
                            height: "100%",
                            backgroundColor: "button.mouse.main",
                            "&:hover": {
                                backgroundColor: "button.mouse.hover",
                            },
                            color: "button.mouse.text",
                        }}
                        onClick={async () => {
                            const [mousePrCon, keyboardPrCon] = preConfig[value]
                            await fetch(`/api/set/mouse?key=${"CLEAR_ALL"}&value=NONE`)
                            for (let key in mousePrCon) {
                                await fetch(`/api/set/mouse?key=${key}&value=${mousePrCon[key]}`)
                            }
                            await fetch(`/api/set/keyboard?key=${"CLEAR_ALL"}&value=NONE`)
                            for (let key in keyboardPrCon) {
                                await fetch(`/api/set/keyboard?key=${key}&value=${keyboardPrCon[key]}`)
                            }   
                        }}
                    >
                        <Typography variant="h6" component="h1" noWrap={true}>
                            {value}
                        </Typography>
                    </Button>
                </Grid>
            ))}
            


        </Grid>
    </div>
}


export default PreConfigPanel;