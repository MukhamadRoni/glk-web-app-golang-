(function () {
    const API_BANK = "/api/v1/bank-soal";
    const API_JP = "/api/v1/jenis-pendidikan";
    const API_MAPEL = "/api/v1/mapel";

    const pathParts = window.location.pathname.split('/');
    const bankId = pathParts[pathParts.length - 1] === "form" ? null : pathParts[pathParts.length - 1];

    const selectJP = document.getElementById("bankJP");
    const selectMapel = document.getElementById("bankMapel");
    const questionsContainer = document.getElementById("questionsContainer");
    
    let questionCount = 0;

    async function loadDropdowns() {
        try {
            const [resJP, resMapel] = await Promise.all([
                fetch(API_JP + "?_t=" + new Date().getTime(), { credentials: "same-origin", cache: "no-store" }),
                fetch(API_MAPEL + "?_t=" + new Date().getTime(), { credentials: "same-origin", cache: "no-store" })
            ]);
            const jpData = await resJP.json();
            const mapelData = await resMapel.json();

            if (jpData && jpData.data) {
                jpData.data.forEach(jp => {
                    const opt = document.createElement('option');
                    opt.value = jp.id;
                    opt.text = `${jp.name} (${jp.jenis_pendidikan})`;
                    selectJP.appendChild(opt);
                });
            }

            if (mapelData && mapelData.data) {
                mapelData.data.forEach(m => {
                    const opt = document.createElement('option');
                    opt.value = m.id;
                    opt.text = `${m.name} (${m.code})`;
                    // If we want to filter Mapel by JP later, we can add data-jpid attribute
                    opt.setAttribute('data-jpid', m.jenis_pendidikan_id);
                    selectMapel.appendChild(opt);
                });
            }
        } catch (error) {
            console.error("Error loading dropdowns:", error);
        }
    }

    function addQuestion(data = null) {
        questionCount++;
        const template = document.getElementById("questionTemplate").content.cloneNode(true);
        const card = template.querySelector('.question-card');
        const qNumber = template.querySelector('.question-number');
        const qText = template.querySelector('.question-text');
        const qType = template.querySelector('.question-type');
        const optContainer = template.querySelector('.options-container');
        const btnAddOpt = template.querySelector('.btn-add-option');
        
        // Setup identifiers
        const qId = `q_${new Date().getTime()}_${questionCount}`;
        card.setAttribute('data-qid', qId);
        
        if (data) {
            qText.value = data.question_text;
            qType.value = data.question_type;
        }

        // Handle question type change
        qType.addEventListener('change', function() {
            renderOptionsUI(this.value, optContainer, qId, btnAddOpt);
        });

        // Handle add option click
        btnAddOpt.addEventListener('click', function() {
            addOptionRow(qType.value, optContainer, qId);
        });

        // Handle remove question
        template.querySelector('.btn-remove-question').addEventListener('click', function() {
            card.remove();
            updateQuestionNumbers();
        });

        questionsContainer.appendChild(template);
        updateQuestionNumbers();

        // Initial render of options
        if (data && data.options && data.options.length > 0) {
            renderOptionsUI(qType.value, optContainer, qId, btnAddOpt, data.options);
        } else {
            renderOptionsUI(qType.value, optContainer, qId, btnAddOpt);
        }
    }

    function updateQuestionNumbers() {
        const cards = questionsContainer.querySelectorAll('.question-card');
        cards.forEach((card, index) => {
            card.querySelector('.question-number').textContent = `Pertanyaan #${index + 1}`;
        });
    }

    function renderOptionsUI(type, container, qId, btnAddOpt, existingOptions = null) {
        container.innerHTML = '';
        if (type === 'MULTIPLE_CHOICE') {
            btnAddOpt.style.display = 'inline-block';
            if (existingOptions) {
                existingOptions.forEach(opt => addOptionRow(type, container, qId, opt));
            } else {
                addOptionRow(type, container, qId);
                addOptionRow(type, container, qId);
            }
        } else {
            btnAddOpt.style.display = 'none';
            if (existingOptions && existingOptions.length > 0) {
                addOptionRow(type, container, qId, existingOptions[0]);
            } else {
                addOptionRow(type, container, qId);
            }
        }
    }

    function addOptionRow(type, container, qId, data = null) {
        let template;
        if (type === 'MULTIPLE_CHOICE') {
            template = document.getElementById("optionTemplatePG").content.cloneNode(true);
            const radio = template.querySelector('.option-correct');
            radio.name = `correct_${qId}`; // Group radios per question
            
            if (data) {
                template.querySelector('.option-text').value = data.option_text;
                if (data.is_correct === 'T') radio.checked = true;
            }

            template.querySelector('.btn-remove-option').addEventListener('click', function(e) {
                e.target.closest('.option-row').remove();
            });
        } else {
            template = document.getElementById("optionTemplateEssay").content.cloneNode(true);
            if (data) {
                template.querySelector('.option-text').value = data.option_text;
            }
        }
        container.appendChild(template);
    }

    async function loadData() {
        if (!bankId) {
            addQuestion(); // Start with one empty question
            return;
        }

        try {
            document.getElementById("formTitle").textContent = "Edit Bank Soal";
            document.getElementById("btnSaveNewVersion").classList.remove("d-none");

            const res = await fetch(`${API_BANK}/${bankId}?_t=` + new Date().getTime(), { credentials: "same-origin", cache: "no-store" });
            const json = await res.json();

            if (json.data) {
                const b = json.data;
                document.getElementById("bankTitle").value = b.title;
                document.getElementById("bankJP").value = b.jenis_pendidikan_id;
                document.getElementById("bankMapel").value = b.mata_pelajaran_id;
                document.getElementById("bankActive").checked = (b.active === 'T');

                if (b.questions && b.questions.length > 0) {
                    b.questions.forEach(q => addQuestion(q));
                } else {
                    addQuestion();
                }
            }
        } catch (error) {
            console.error("Error loading bank soal detail:", error);
            Swal.fire('Error', 'Gagal memuat data.', 'error');
        }
    }

    document.getElementById("btnAddQuestion").addEventListener('click', () => addQuestion());

    document.getElementById("formBankSoal").addEventListener('submit', async function(e) {
        e.preventDefault();
        saveBankSoal(false);
    });

    document.getElementById("btnSaveNewVersion").addEventListener('click', async function(e) {
        e.preventDefault();
        // Trigger HTML5 validation check
        if (!document.getElementById("formBankSoal").checkValidity()) {
            document.getElementById("formBankSoal").reportValidity();
            return;
        }
        saveBankSoal(true);
    });

    async function saveBankSoal(isNewVersion) {
        const payload = {
            title: document.getElementById("bankTitle").value,
            jenis_pendidikan_id: parseInt(document.getElementById("bankJP").value),
            mata_pelajaran_id: parseInt(document.getElementById("bankMapel").value),
            active: document.getElementById("bankActive").checked ? "T" : "F",
            save_as_new_version: isNewVersion,
            questions: []
        };

        const cards = questionsContainer.querySelectorAll('.question-card');
        let hasError = false;

        cards.forEach((card, index) => {
            const qType = card.querySelector('.question-type').value;
            const qText = card.querySelector('.question-text').value;

            const qPayload = {
                question_type: qType,
                question_text: qText,
                order_index: index,
                options: []
            };

            const optionRows = card.querySelectorAll('.option-row');
            if (qType === 'MULTIPLE_CHOICE') {
                let hasCorrect = false;
                optionRows.forEach(row => {
                    const text = row.querySelector('.option-text').value;
                    const isCorrect = row.querySelector('.option-correct').checked ? "T" : "F";
                    if (isCorrect === 'T') hasCorrect = true;
                    qPayload.options.push({
                        option_text: text,
                        is_correct: isCorrect
                    });
                });
                if (!hasCorrect) {
                    Swal.fire('Peringatan', `Pertanyaan #${index + 1} harus memiliki minimal 1 jawaban benar!`, 'warning');
                    hasError = true;
                }
            } else {
                optionRows.forEach(row => {
                    qPayload.options.push({
                        option_text: row.querySelector('.option-text').value,
                        is_correct: "T" // For essay, just default to T
                    });
                });
            }

            payload.questions.push(qPayload);
        });

        if (hasError) return;

        try {
            const url = bankId ? `${API_BANK}/${bankId}` : API_BANK;
            const method = bankId ? "PUT" : "POST";

            const res = await fetch(url, {
                method: method,
                headers: { "Content-Type": "application/json" },
                credentials: "same-origin",
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                await Swal.fire('Berhasil!', 'Bank Soal berhasil disimpan.', 'success');
                window.location.href = "/admin/recruitment/master/bank-soal";
            } else {
                const errData = await res.json();
                Swal.fire('Gagal', "Gagal menyimpan: " + (errData.message || errData.error), 'error');
            }
        } catch (error) {
            console.error(error);
            Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
        }
    }

    // Init
    loadDropdowns().then(() => loadData());
})();
