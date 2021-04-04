const log4js = require('log4js');
const logger = log4js.getLogger('Handlers');
const fs = require('fs');
const yaml = require('js-yaml');
const {FileSystemWallet, X509WalletMixin, Gateway } = require('fabric-network');
logger.level = 'debug';
const path = require('path');
const fixtures = path.resolve(__dirname, '../');
const configfilepath = '../connection.yaml';

var contract;
async function getcontract() {
    contract = await identity();
}

getcontract();

async function identity() {
    const gateway = new Gateway();
    const wallet = await new FileSystemWallet('../identity/user/isabella/wallet');

    // Main try/catch block
    try {

	const credPath = path.join(fixtures, '/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com');
        const cert = fs.readFileSync(path.join(credPath, '/msp/signcerts/User1@org1.example.com-cert.pem')).toString();
        const key = fs.readFileSync(path.join(credPath, '/msp/keystore/c75bd6911aca808941c3557ee7c97e90f3952e379497dc55eb903f31b50abc83_sk')).toString();

        // Load credentials into wallet
        const identityLabel = 'User1@org1.example.com';
        const identity = X509WalletMixin.createIdentity('Org1MSP', cert, key);

        await wallet.import(identityLabel, identity);


        let connectionProfile = yaml.safeLoad(fs.readFileSync(configfilepath, 'utf8'));

        logger.debug('==========================================');
        //Enroll the admin user, and import the new identity into the wallet.
        let connectionOptions = {
            identity: identityLabel,
            wallet: wallet,
            discovery: { enabled: false, asLocalhost: true }

        };

       await gateway.connect(connectionProfile, connectionOptions);
        logger.debug('Use network channel: mychannel.');

        logger.debug('Gateway connected');

        const network = await gateway.getNetwork("mychannel");
        logger.debug('Got network');
        const contract = await network.getContract('mycc');
        logger.debug('Contract set and returned');
        return contract;
    } catch (error) {
        logger.debug(`Error processing transaction. ${error}`);
        logger.debug(error.stack);
        return null;
    }
}


const addPatient = async (req, res) => {
    logger.debug('===================POST API /patient called==================');
    var object = req.body.object;
    try {
        const result = await contract.submitTransaction('addPatient', object);
        return res.status(200).json({msg: result.toString()});

    } catch (err) {
        logger.error("Error received ", err)
        return res.status(500).json(response.error(err));
    } finally {
        logger.debug('Transaction ends')
    }
};


const getPatient = async (req, res) => {
    logger.debug('===================GET API /patient/:id called==================');
    var patientid = req.params.id;
    try {
        const result = await contract.evaluateTransaction('getPatient', patientid);
        let object = response.res("Received the patient details", result.toString());
        return res.status(200).json(object);
    } catch (err) {
        logger.error("Error received ", err)
        return res.status(500).json(response.error(err));
    }
};


const editPatient = async (req, res) => {
    logger.debug('===================PUT API /patient called==================');
    var object = req.body.object;
    var userId = req.params.id;

    try {
        const result = await contract.submitTransaction('editPatient', object, userId);
        return res.status(200).json({msg: result.toString()});
    } catch (err) {
        logger.error("Error received ", err)
        return res.status(500).json(response.error(err));
    } finally {
        logger.debug('Transaction ends')
    }
};

exports.addPatient = addPatient;
exports.getPatient = getPatient;
exports.editPatient = editPatient;
