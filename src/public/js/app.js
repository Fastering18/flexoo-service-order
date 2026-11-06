document.addEventListener("DOMContentLoaded", function () {
    var userTrigger = document.getElementById("userMenuTrigger");
    var userDropdown = document.getElementById("userDropdown");
    if (userTrigger && userDropdown) {
        userTrigger.addEventListener("click", function (e) {
            e.stopPropagation();
            userDropdown.classList.toggle("open");
        });
        document.addEventListener("click", function () {
            userDropdown.classList.remove("open");
        });
    }

    var scaleBtns = document.querySelectorAll(".scale-btn");
    var scaleInput = document.getElementById("scaleInput");
    var supportNote = document.getElementById("supportNote");
    scaleBtns.forEach(function (btn) {
        btn.addEventListener("click", function () {
            scaleBtns.forEach(function (b) { b.classList.remove("active"); });
            btn.classList.add("active");
            var val = btn.getAttribute("data-scale");
            if (scaleInput) scaleInput.value = val;
            if (supportNote) supportNote.style.display = val === "Enterprise (custom)" ? "block" : "none";
        });
    });

    var dropzone = document.getElementById("srsDropzone");
    var fileInput = document.getElementById("srsFileInput");
    var fileNameInput = document.getElementById("srsFileNameInput");
    var selectedFileLabel = document.getElementById("selectedFileName");
    function setDemoFile(name) {
        if (fileNameInput) fileNameInput.value = name;
        if (selectedFileLabel) selectedFileLabel.textContent = "Selected (Demo No-Op): " + name;
    }
    if (dropzone && fileInput) {
        dropzone.addEventListener("click", function () { fileInput.click(); });
        fileInput.addEventListener("change", function () {
            if (fileInput.files && fileInput.files[0]) {
                setDemoFile(fileInput.files[0].name);
                fileInput.value = "";
            }
        });
        dropzone.addEventListener("dragover", function (e) { e.preventDefault(); });
        dropzone.addEventListener("drop", function (e) {
            e.preventDefault();
            if (e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files[0]) {
                setDemoFile(e.dataTransfer.files[0].name);
            }
        });
    }

    var modalBackdrop = document.getElementById("orderModalBackdrop");
    var modalClose = document.getElementById("orderModalClose");
    document.querySelectorAll(".order-card").forEach(function (card) {
        card.addEventListener("click", function () {
            var id = card.getAttribute("data-order-id");
            document.getElementById("modalOrderId").textContent = id;
            document.getElementById("modalProjectName").textContent = card.getAttribute("data-project-name");
            document.getElementById("modalServiceName").textContent = card.getAttribute("data-service-name");
            document.getElementById("modalScale").textContent = card.getAttribute("data-scale");
            document.getElementById("modalStatus").textContent = card.getAttribute("data-status");
            document.getElementById("modalSrs").textContent = card.getAttribute("data-srs") || "N/A";
            document.getElementById("modalDelivery").textContent = card.getAttribute("data-delivery");
            document.getElementById("modalComments").textContent = card.getAttribute("data-comments") || "-";
            document.getElementById("modalStatusForm").action = "/orders/" + encodeURIComponent(id) + "/status";
            document.getElementById("modalStatusSelect").value = card.getAttribute("data-status");
            modalBackdrop.classList.add("open");
        });
    });
    if (modalClose) modalClose.addEventListener("click", function () { modalBackdrop.classList.remove("open"); });
    if (modalBackdrop) {
        modalBackdrop.addEventListener("click", function (e) {
            if (e.target === modalBackdrop) modalBackdrop.classList.remove("open");
        });
    }
});

