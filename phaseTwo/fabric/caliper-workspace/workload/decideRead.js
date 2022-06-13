'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const { ENOBUFS } = require('constants');
require('dotenv').config();

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        var cliattr = {
            "status": true, "expiration": "02-Jan-2026",
            "organisationid": "spain01"
        };
        this.clientAttr = JSON.stringify(cliattr);
        this.randomIds = [];
    }


    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        var attr_cli = {
            "status": true, "expiration": "02-Jan-2026",
            "organisationid": "spain01"
        };
        const desID = `${this.workerIndex}_${'idxxx'}`;
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DecideRead',
            invokerIdentity: 'User1',
            contractArguments: [desID, 'did:client:123456789abcdefghigklmn', 'r004', JSON.stringify(attr_cli)],
            readOnly: false,
            timeout: 180,
        };
        await this.sutAdapter.sendRequests(request);
    }

    async submitTransaction() {
        var randomId = Math.floor(Math.random() * this.roundArguments.decisions);
        while (this.randomIds.includes(randomId)) {
            randomId = Math.floor(Math.random() * this.roundArguments.decisions);
        }
        this.randomIds.push(randomId);
        const desID = `${this.workerIndex}_${randomId}`;
        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DecideRead',
            invokerIdentity: 'User1',
            contractArguments: [desID, 'did:client:123456789abcdefghigklmn_' + randomId, 'r004', this.clientAttr],
            readOnly: false,
            timeout: 180,
        };
        await this.sutAdapter.sendRequests(myArgs);
    }

    async cleanupWorkloadModule() {
        for (let i = 0; i < this.roundArguments.decisions; i++) {
            if (this.randomIds.includes(i)) {
                const desID = `${this.workerIndex}_${i}`;
                const request = {
                    contractId: this.roundArguments.contractId,
                    contractFunction: 'DeleteDecision',
                    invokerIdentity: 'User1',
                    contractArguments: [desID],
                    readOnly: false
                };
                await this.sutAdapter.sendRequests(request);
            }
        }
        const desID = `${this.workerIndex}_${'idxxx'}`;
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DeleteDecision',
            invokerIdentity: 'User1',
            contractArguments: [desID],
            readOnly: false
        };
        await this.sutAdapter.sendRequests(request);
    }
}


function createWorkloadModule() {
    return new MyWorkload();
}


module.exports.createWorkloadModule = createWorkloadModule;
