// Copyright (C) 2021-2022 David Sugar <tychosoft@gmail.com>.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package osip

type SIP_STATUS int
type EVT_TYPE string

const (
	SIP_TRYING     SIP_STATUS = 100
	SIP_RINGING    SIP_STATUS = 180
	SIP_FORWARDING SIP_STATUS = 181
	SIP_QUEUED     SIP_STATUS = 182
	SIP_PROGRESS   SIP_STATUS = 183
	SIP_TERMINATED SIP_STATUS = 199

	SIP_OK              SIP_STATUS = 200
	SIP_ACCEPTED        SIP_STATUS = 202
	SIP_NO_NOTIFICATION SIP_STATUS = 204

	SIP_MULTIPLE_CHOICES    SIP_STATUS = 300
	SIP_MOVED_PERMENANTLY   SIP_STATUS = 301
	SIP_MOVED_TEMPORARILY   SIP_STATUS = 302
	SIP_USE_PROXY           SIP_STATUS = 305
	SIP_ALTERNATIVE_SERVICE SIP_STATUS = 380

	SIP_BAD_REQUEST              SIP_STATUS = 400
	SIP_UNAUTHORIZED             SIP_STATUS = 401
	SIP_PAYMENT_REQUIRED         SIP_STATUS = 402
	SIP_FORBIDDEN                SIP_STATUS = 403
	SIP_NOT_FOUND                SIP_STATUS = 404
	SIP_METHOD_NOT_ALLOWED       SIP_STATUS = 405
	SIP_NOT_ACCEPTABLE           SIP_STATUS = 406
	SIP_PROXY_AUTH_REQUIRED      SIP_STATUS = 407
	SIP_REQUEST_TIMEOUT          SIP_STATUS = 408
	SIP_CONFLICT                 SIP_STATUS = 409
	SIP_GONE                     SIP_STATUS = 410
	SIP_LENGTH_REQUIRED          SIP_STATUS = 411
	SIP_CONDITIONAL_FAILED       SIP_STATUS = 412
	SIP_REQUEST_TOO_LARGE        SIP_STATUS = 413
	SIP_URI_TOO_LONG             SIP_STATUS = 414
	SIP_UNSUPPORTED_MEDIA        SIP_STATUS = 415
	SIP_UNSUPPORTED_URI          SIP_STATUS = 416
	SIP_UNKNOWN_PRIORITY         SIP_STATUS = 417
	SIP_BAD_EXTENSION            SIP_STATUS = 420
	SIP_EXTENSION_REQUIRED       SIP_STATUS = 421
	SIP_INVALID_SESSION_INTERVAL SIP_STATUS = 422
	SIP_INTERVAL_TOO_BRIEF       SIP_STATUS = 423
	SIP_BAD_LOCATIOPN            SIP_STATUS = 424
	SIP_BAD_ALERT                SIP_STATUS = 425
	SIP_MISSING_IDENTITY         SIP_STATUS = 428
	SIP_MISSING_REFERRER         SIP_STATUS = 429
	SIP_FLOW_FAILED              SIP_STATUS = 430
	SIP_ANONYMITY_DISALLOWED     SIP_STATUS = 433
	SIP_BAD_IDENTITY             SIP_STATUS = 436
	SIP_INVALID_CERT             SIP_STATUS = 437
	SIP_INVALID_IDENTITY         SIP_STATUS = 438
	SIP_FIRST_HOP_NO_OUTBOUND    SIP_STATUS = 439
	SIP_MAX_BREADTH_EXCEEDED     SIP_STATUS = 440
	SIP_BAD_INFO                 SIP_STATUS = 469
	SIP_CONSENT_NEEDED           SIP_STATUS = 470
	SIP_TEMPORARILY_UNAVAILABLE  SIP_STATUS = 480
	SIP_CALL_DOES_NOT_EXIST      SIP_STATUS = 481
	SIP_LOOP_DETECT              SIP_STATUS = 482
	SIP_TOO_MANY_HOPS            SIP_STATUS = 483
	SIP_ADDRESS_INCOMPLETE       SIP_STATUS = 484
	SIP_AMBIGUOUS                SIP_STATUS = 485
	SIP_BUSY_HERE                SIP_STATUS = 486
	SIP_REQUEST_TERMINATED       SIP_STATUS = 487
	SIP_NOT_ACCEPTABLE_HERE      SIP_STATUS = 488
	SIP_BAD_EVENT                SIP_STATUS = 489
	SIP_REQUEST_PENDING          SIP_STATUS = 491
	SIP_UNDECIPHERABLE           SIP_STATUS = 493
	SIP_SECURITY_REQUIRED        SIP_STATUS = 494

	SIP_INTERNAL_ERROR          SIP_STATUS = 500
	SIP_NOT_IMPLIMENTED         SIP_STATUS = 501
	SIP_BAD_GATEWAY             SIP_STATUS = 502
	SIP_SERVICE_UNAVAILABLE     SIP_STATUS = 503
	SIP_SERVER_TIMEOUT          SIP_STATUS = 504
	SIP_VERSION_UNSUPPORTED     SIP_STATUS = 505
	SIP_MESSAGE_TOO_LARGE       SIP_STATUS = 513
	SIP_PUSH_NOTIFY_UNSUPPORTED SIP_STATUS = 555
	SIP_PPRECONDITION_FAILED    SIP_STATUS = 580

	SIP_BUSY_EVERYWHERE         SIP_STATUS = 600
	SIP_DECLINE                 SIP_STATUS = 603
	SIP_DOES_NOT_EXIST_ANYWHERE SIP_STATUS = 604
	SIP_NOT_ACCEPTABLE_ANYWHERE SIP_STATUS = 606
	SIP_UNWANTED                SIP_STATUS = 607
	SIP_REJECTED                SIP_STATUS = 608
)

const (
	EVT_TIMEOUT  EVT_TYPE = "timeout"
	EVT_STARTUP  EVT_TYPE = "startup"
	EVT_SHUTDOWN EVT_TYPE = "shutdown"
)
