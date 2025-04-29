
/* merged from 2024-07-13 */

function checkEmailPattern(userEmail) {
    let atPos = 0;
    let atCount = 0;
    for (let i = 0; i < userEmail.length; i++) {
        if (userEmail[i] === '@') {
            atPos = i;
            atCount++;
        }
    }
    if (atCount !== 1 || atPos === 0 || atPos + 1 === userEmail.length) {
        return false;
    }
    let dotCount = 0;
    for (let i = atPos + 1; i < userEmail.length; i++) {
        if (userEmail[i] === '.') {
            dotCount++;
        }
    }
    if (dotCount === 0 || userEmail[atPos + 1] === '.' || userEmail[userEmail.length - 1] === '.') {
        return false;
    }
    return true;
}

const hexBase = 16;
const hexDigits = "0123456789abcdef";
function toHex(number) {
    let res = "";
    while (number > 0) {
        res += hexDigits[number % BigInt(hexBase)];
        number = BigInt(number / BigInt(hexBase));
    }
    return res;
}

const p1 = 307; 
const mod1 = 1000000007;
const p2 = 311;
const mod2 = 1000000009;
const p3 = 313;
const mod3 = 1000000021;
const p4 = 317;
const mod4 = 1000000033;

const base1 = 1000000093;
const base2 = 1000000097;
const base3 = 1000000103;
const base4 = 1000000123;

function getPasswordHash(pass) {
    let hash1 = BigInt(0);
    let hash2 = BigInt(0);
    let hash3 = BigInt(0);
    let hash4 = BigInt(0);
    for (let i = 0; i < pass.length; i++) {
        let char = BigInt(pass.charCodeAt(i));
        hash1 = BigInt(BigInt(hash1 * BigInt(p1)) % BigInt(mod1) + char) % BigInt(mod1);
        hash2 = BigInt(BigInt(hash2 * BigInt(p2)) % BigInt(mod2) + char) % BigInt(mod2);
        hash3 = BigInt(BigInt(hash3 * BigInt(p3)) % BigInt(mod3) + char) % BigInt(mod3);
        hash4 = BigInt(BigInt(hash4 * BigInt(p4)) % BigInt(mod4) + char) % BigInt(mod4);
    }
    hash1 += BigInt(base1);
    hash2 += BigInt(base2);
    hash3 += BigInt(base3);
    hash4 += BigInt(base4);
    return toHex(hash1) + toHex(hash2) + toHex(hash3) + toHex(hash4);
}

function checkFormValues() {
    userName = document.getElementById("username").value;
    userEmail = document.getElementById("email").value;
    passwordOne = document.getElementById("pass_one").value;
    passwordTwo = document.getElementById("pass_two").value;
    if (userName.length < 4) {
        alert("Username should be at least with 4 characaters");
        return;
    }
    if (!checkEmailPattern(userEmail)) {
        alert("Invalid email. Possible email pattern:\n youremail@box.domen");
        return;
    }
    const hash1 = getPasswordHash(passwordOne);
    const hash2 = getPasswordHash(passwordTwo);
    if (hash1 != hash2) {
        alert("Two passwords are not the same");
        return;
    }
    const NeededHand = "http://localhost:8680/api/credentials/?username="+userName+"&useremail="+userEmail; 
    fetch(NeededHand).then((response) => response.json()).then((json) => {
        console.log(json);
        console.log(json.UserEmailResponse);
        console.log(json.UserNameResponse);                              
        if (json.UserNameResponse === "EXISTS") {
            alert("This username is alredy utilized!");
            return;
        }
        if ((String(json.UserNameResponse) === "OK") === false) {
            alert("An error occured " + json.UserNameResponse);
            return;
        } 
        if (json.UserEmailResponse === "EXISTS") {
            alert("This email is already utilized!");
            return;
        }
        if ((String(json.UserEmailResponse) === "OK") === false) {
            alert("An error occured " + json.UserEmailResponse);
            return;
        }
        document.cookie = "password=" + hash1 + "; path=/loader/";
        const hiddenFormSender = document.createElement('form');
        hiddenFormSender.method = 'POST';
        hiddenFormSender.action = '/loader/?action=register';
        const senderParams = {
            NameOfUser: userName,
            EmailOfUser: userEmail
        }
        for (const key in senderParams) {
            if (senderParams.hasOwnProperty(key)) {
                const hiddenInput = document.createElement('input');
                hiddenInput.type = 'text';
                hiddenInput.name = key;
                hiddenInput.value = senderParams[key];
                hiddenInput.classList.add("ghost_box");
                hiddenFormSender.appendChild(hiddenInput);
            }
        }
        document.body.appendChild(hiddenFormSender);
        hiddenFormSender.submit();
    });
}

function checkLoginValues() {
    let userEmail = document.getElementById("email").value;
    let userPass = document.getElementById("pass_one").value;
    userPass = getPasswordHash(userPass);
    if (!checkEmailPattern(userEmail)) {
        alert("Invalid email. Possible email pattern:\n youremail@box.domen");
        return;
    }
    document.cookie = "password=" + userPass + "; path=/loader/";
    const hiddenFormSender = document.createElement('form');
    hiddenFormSender.method = 'POST';
    hiddenFormSender.action = '/loader/?action=login';
    const hiddenInput = document.createElement('input');
    hiddenInput.type = 'text';
    hiddenInput.name = 'EmailOfUser';
    hiddenInput.value = userEmail;
    hiddenInput.classList.add("ghost_box");
    hiddenFormSender.appendChild(hiddenInput);
    document.body.appendChild(hiddenFormSender);
    hiddenFormSender.submit();
}
