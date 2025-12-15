$(document).ready(() => {
    fill_information();
});

function fill_information() {
    fill_instructions();
    fill_state();
    fill_memory();
}

function fill_instructions() {
    $.get("/instructions", (data) => {
        var instruction_div = "";
        current_pc = data["Pc"];

        for (const [address, instruction] of Object.entries(data["Instructions"])) {
            // Compare using instruction struct Pc field to be safe, loose equality for string/number match
            let isCurrent = (instruction["Pc"] == current_pc);
            let div_class = isCurrent ? "disassembly-line current" : "disassembly-line next";

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
        fill_information();
    });
}

function fill_state() {
    $.get("/cpu-state", (data) => {
        $("#cpu-state #pc").html(data["PC"].toString(16).toUpperCase());
        $("#cpu-state #a").html(data["A"].toString(16).toUpperCase());
        $("#cpu-state #x").html(data["X"].toString(16).toUpperCase());
        $("#cpu-state #y").html(data["Y"].toString(16).toUpperCase());
        $("#cpu-state #sp").html(data["SP"].toString(16).toUpperCase());

        let flags = data["Flags"];

        $("#cpu-state #carry").toggleClass("active", flags["Carry"]);
        $("#cpu-state #zero").toggleClass("active", flags["Zero"]);
        $("#cpu-state #interrupt-disable").toggleClass("active", flags["InterruptDisable"]);
        $("#cpu-state #decimal-mode").toggleClass("active", flags["DecimalMode"]);
        $("#cpu-state #b").toggleClass("active", flags["B"]);
        $("#cpu-state #overflow").toggleClass("active", flags["Overflow"]);
        $("#cpu-state #negative").toggleClass("active", flags["Negative"]);
    });
}

function dump_to_string(dump) {
    var str = "";
    for (const [address, value] of Object.entries(dump)) {
        let address_string = ("0000" + parseInt((address)).toString(16).toUpperCase()).slice(-4);
        str += `[${address_string}] ${value}<br>`;
    }
    return str;
}

function fill_memory() {
    $.get("/memory-dump", (data) => {
        var zero_page = dump_to_string(data["ZeroPage"]);
        $("#zero-page-dump").html(zero_page);
        var stack = dump_to_string(data["Stack"]);
        $("#stack-dump").html(stack);
    });
}