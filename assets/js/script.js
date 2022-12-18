// Aulia Pradina, 1IA26, 50421237

// Activate PopOver
var popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'));
var popoverList = popoverTriggerList.map(function(popoverTrigger) {
    return new bootstrap.Popover(popoverTrigger);
});