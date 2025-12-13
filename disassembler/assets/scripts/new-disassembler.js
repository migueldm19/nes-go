$(document).ready(() => {
    fill_instructions();
})

function fill_instructions() {
    $.get("/instructions", (data) => {
        var instruction_div = "";
        current_pc = data["Pc"];

        for (const [address, instruction] of Object.entries(data["Instructions"])) {
            let div_class = (address == current_pc) ? "current" : "next";
            instruction_div +=
                '<div class="' + div_class + '">' +
                    '<span class="address">' +
                        instruction["Pc"].toString(16).toUpperCase() +
                    '</span> ' +
                    instruction["InstructionText"] +
                '</div>';
        }

        $("#instructions").html(instruction_div);

        // Scroll to the current instruction
        var currentElement = $("#instructions .current");
        if (currentElement.length) {
            currentElement[0].scrollIntoView({
                behavior: 'auto',
                block: 'center',  // or 'nearest', 'start', 'center', 'end'
                inline: 'nearest'
            });
        }
    });
}

function step_disassembler() {
    $.post("/step", (data) => {
        fill_instructions()
    });
}

function fill_state() {
    $.get("/cpu-state", (data) => {

    });
}