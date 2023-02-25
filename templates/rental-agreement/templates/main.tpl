{{define "checklist-row"}}
   {{$margin := getWidth (printf "%s%s" . " Condition ") "Times-Roman"}}
   <table-cell>
      <division margin="5 0">
         <paragraph>
            <text-chunk font = "times">{{.}} Condition </text-chunk>
         </paragraph>
         <line fit-mode="fill-width" position="relative" thickness= "0.5" margin="0 0 0 {{$margin}}"></line>
      </division>
   </table-cell>
   <table-cell>
      <division margin="5 0">
         <paragraph>
            <text-chunk>Specific Damage </text-chunk>
         </paragraph>
         <line fit-mode="fill-width" position="relative" thickness= "0.5" margin="0 0 0 80"></line>
      </division>
   </table-cell>
{{end}}

{{define "check-list-table"}}
   <table columns="2" margin="20 0 0 0">
      {{range .Items}}
         {{template "checklist-row" .}}
      {{end}}
   </table>
{{end}}

{{define "form-sig"}}
   <division>
      <paragraph>
         <text-chunk font= "times-bold" font-size="11">{{.Text}}: </text-chunk>
      </paragraph>
      <line fit-mode="fill-width" position="relative" thickness= "0.2" margin="{{.Margin}}"></line>
   </division>
{{end}}

{{define "simple-form"}}
      <paragraph>
         <text-chunk font= "times" font-size="11">{{.Text}} </text-chunk>
      </paragraph>
      <line fit-mode="fill-width" position="relative" thickness= "0.2" margin="{{.Margin}}"></line>
{{end}}

{{define "paragraph-with-header"}}
   {{$margin := "18 0 0 0"}}
   {{if .Margin}}
      {{$margin:= .Margin}}
   {{end}}
   <paragraph margin= "{{$margin}}" line-height="1.1">
      <text-chunk font="times-bold" font-size="12">{{.Header}}: </text-chunk>
      <text-chunk font="times" font-size="12">{{.Text}} </text-chunk>
   </paragraph>
{{end}}

{{define "simple-paragraph"}}
   {{$fontSize := 12}}
   {{$margin := "18 0 0 0"}}
   {{$font := "times"}}
   {{if .FontSize}}
      {{$fontSize = .FontSize}}
   {{end}}
   {{if .Margin}}
      {{$margin = .Margin}}
   {{end}}
   {{if .Font}}
      {{$font = .Font}}
   {{end}}
   <paragraph margin="{{$margin}}" line-height="1.1">
      <text-chunk font="{{$font}}" font-size="{{$fontSize}}">{{.Text}}</text-chunk>
   </paragraph>
{{end}}

<paragraph margin="0 0 10 0" text-align="center" line-height="1.1">
   <text-chunk font="times-bold" font-size="20"> LEASE WITH OPTION TO PURCHASE </text-chunk>
</paragraph>

{{template "simple-paragraph" dict "Text" (printf `This agreement, dated %s, by and between a business entity known as %s of %s, hereinafter known as the “Landlord”.` (formatTime .Date "December 9 2006") .Company.Name .Company.Address)}}

{{template "simple-paragraph" dict "Margin" "16 0 0 0" "Font" "times-bold" "Text" "AND"}}

{{template "simple-paragraph" dict "Text" (printf "%d individuals known as %s, hereinafter known as the “Tenant(s)”, agree to the following:" (len .Tenant.Names) (listItems .Tenant.Names true))}}

{{template "paragraph-with-header" dict "Header" "OCCUPANT(S)" "Text" (printf `The Premises is to be occupied strictly as a residential dwelling with the
following %s (%d) Occupants to reside on the Premises in addition to the Tenant(s) mentioned above: %s, hereinafter known as the “Occupant(s)”.` (numberToWord (len .Tenant.Names) true) (len .Tenant.Names) (listItems .Tenant.Names true))}}

{{template "paragraph-with-header" dict "Header" "OFFER TO RENT" "Text" (printf `The Landlord hereby rents to the Tenant(s), subject to the following terms 
and conditions of this Agreement, an apartment with the address of %s consisting of %.1f bathroom(s) and %d bedroom(s) hereinafter known as 
the “Premises”. The Landlord may also use the address for notices sent to the Tenant(s).` .Apartment.Address .Apartment.Bathrooms .Apartment.Bedrooms)}}

{{template "paragraph-with-header" dict "Header" "PURPOSE" "Text" `The Tenant(s) and any Occupant(s) may only use the Premises as a residential 
dwelling. It may not be used for storage, manufacturing of any type of food or product, 
professional service(s), or for any commercial use unless otherwise stated in this Agreement.`}}

{{template "paragraph-with-header" dict "Header" "FURNISHINGS" "Text" "The Premises is furnished with the following:"}}

{{template "simple-paragraph" dict "Text" (printf `%s and all other furnishings to be 
provided by the Tenant(s). Any damage to the Landlord's furnishings shall be the liability of the Tenant(s), reasonable wear-and-tear excepted, to be billed directly or less the Security Deposit.` (listItems .Apartment.FurnishingItems false))}}

{{template "paragraph-with-header" dict "Header" "APPLIANCES" "Text" "The Landlord shall provide the following appliances:"}}

{{template "simple-paragraph" dict "Margin" "10 0 0 0" "Text" (printf `%s and any
other unnamed appliances existing on the Premises. Any damage to the Landlord's appliances
shall be the liability of the Tenant(s), reasonable wear-and-tear excepted, to be billed directly or less the Security Deposit.` (listItems .Apartment.ProvidedAppliances false))}}

{{template "paragraph-with-header" dict "Header" "LEASE TERM" "Text" (printf `This Agreement shall be a fixed-period arrangement beginning on %s and ending on %s with the Tenant(s) having the option to continue to 
occupy the Premises under the same terms and conditions of this Agreement under a 
Month-to-Month arrangement (Tenancy at Will) with either the Landlord or Tenant having the 
option to cancel the tenancy with at least %s(%d) days notice or the minimum time-period set 
by the State, whichever is shorter. For the Tenant to continue under Month-to-Month tenancy at the expiration of the Lease Term, the Landlord must be notified within %s (%d) days before 
the end of the Lease Term. Hereinafter known as the “Lease Term”.` (formatTime .BeginningDate "December 9 2006") (formatTime .EndingDate "December 9 2006") (numberToWord .CancellationNotificationPeriod false) .CancellationNotificationPeriod (numberToWord .ContinuationNotificationPeriod false) .ContinuationNotificationPeriod)}}

{{template "paragraph-with-header" dict "Header" "RENT" "Text" (printf `Tenant(s) shall pay the Landlord in equal monthly installments of $%s (US
Dollars) hereinafter known as the “Rent”. The Rent will be due on the First (1st) of every
month and be paid through an electronic payment known as Automated Clearing House or 
“ACH”. Details of the Tenant's banking information and authorization shall be attached to this 
Lease Agreement.` .MonthlyInstallment)}}

{{template "paragraph-with-header" dict "Header" "NON-SUFFICIENT FUNDS (NSF CHECKS)" "Text" (printf `If the Tenant(s) attempts to pay the rent with 
a check that is not honored or an electronic transaction (ACH) due to insufficient funds (NSF) 
there shall be a fee of $%s (US Dollars).` .InsufficientFundFee)}}

{{template "paragraph-with-header" dict "Header" "LATE FEE" "Text" `If rent is not paid on the due date, there shall be a late fee assessed by the
Landlord in the amount of: `}}

{{template "simple-paragraph" dict "Text" (printf `$%s (US Dollars) per occurrence for each month payment that is late after the 3rd Day rent
is due.` .LatePaymentFee)}}

{{template "paragraph-with-header" dict "Header" "FIRST (1ST) MONTH'S RENT" "Text" "First (1st) month's rent shall be due by the Tenant(s) upon the execution of this Agreement."}}
{{template "paragraph-with-header" dict "Header" "PRE-PAYMENT" "Text" "The Landlord shall not require any pre-payment of rent by the Tenant(s)."}}
{{template "paragraph-with-header" dict "Header" "PROBATION PERIOD" "Text" "he Tenant(s) will not move into the Premises before the start of the Lease Term."}}

{{template "paragraph-with-header" dict "Header" "SECURITY DEPOSIT" "Text" (printf `A Security Deposit in the amount of $%s (US Dollars) shall be
required by the Tenant(s) at the execution of this Agreement to the Landlord for the faithful
performance of all the terms and conditions. The Security Deposit is to be returned to the
Tenant(s) within %d days after this Agreement has terminated, less any damage charges and
without interest. This Security Deposit shall not be credited towards rent unless the Landlord
gives their written consent.` .SecurityDeposit .SecurityDepositReturnTime)}}

{{template "paragraph-with-header" dict "Header" "POSSESSION" "Text" `Tenant(s) has examined the condition of the Premises and by taking 
possession acknowledges that they have accepted the Premises in good order and in its current 
condition except as herein otherwise stated. Failure of the Landlord to deliver possession of the 
Premises at the start of the Lease Term to the Tenant(s) shall terminate this Agreement at the 
option of the Tenant(s). Furthermore, under such failure to deliver possession by the Landlord, 
and if the Tenant(s) cancels this Agreement, the Security Deposit (if any) shall be returned to 
the Tenant(s) along with any other pre-paid rent, fees, including if the Tenant(s) paid a fee 
during the application process before the execution of this Agreement.`}}


{{template "paragraph-with-header" dict "Header" "OPTION TO PURCHASE" "Text" (printf `The Tenant(s) shall have the right to purchase the Premises
described herein for $%s at any time during the course of the Lease Term, along with
any renewal periods or extensions, by providing written notice to the Landlord along with a
deposit of $%s that is subject to the terms and conditions of a Purchase and Sale
Agreement to be negotiated, in “good faith”, between the Landlord and Tenant(s).` .PurchaseAmount .PurchaseDepositAmount)}}

{{template "simple-paragraph" dict "Text" `If the Landlord and Tenant(s) cannot produce a signed Purchase and Sale Agreement within a
reasonable time period then the deposit shall be refunded to the Tenant(s) and this Lease
Agreement shall continue under its terms and conditions.`}}

{{template "simple-paragraph" dict "Text" `If the option to purchase is exercised by the Tenant(s) all Rent that is paid to the Landlord shall 
remain separate from any and all deposits, consideration, or payments, made to the Landlord in 
regards to the purchase of the Premises.`}}

{{template "paragraph-with-header" dict "Header" "RECORDING" "Text" `The Tenant(s) shall be withheld from recording this Option to Purchase unless 
the Tenant(s) has the written consent from the Landlord.`}}

{{template "paragraph-with-header" dict "Header" "ACCESS" "Text" `Upon the beginning of the Proration Period or the start of the Lease Term, 
whichever is earlier, the Landlord agrees to give access to the Tenant(s) in the form of keys,
fobs, cards, or any type of keyless security entry as needed to enter the common areas and the 
Premises. Duplicate copies of the access provided may only be authorized under the consent of 
the Landlord and, if any replacements are needed, the Landlord may provide them for a fee. At 
the end of this Agreement all access provided to the Tenant(s) shall be returned to the Landlord 
or a fee will be charged to the Tenant(s) or the fee will be subtracted from the Security Deposit.`}}

{{template "paragraph-with-header" dict "Header" "MOVE-IN INSPECTION" "Text" `Before, at the time of the Tenant(s) accepting possession, or
shortly thereafter, the Landlord and Tenant(s) shall perform an inspection documenting the
present condition of all appliances, fixtures, furniture, and any existing damage within the Premises.
`}}

{{template "paragraph-with-header" dict "Header" "SUBLETTING" "Text" `The Tenant(s) shall not have the right to sub-let the Premises or any part
thereof without the prior written consent of the Landlord. If consent is granted by the Landlord,
the Tenant(s) will be responsible for all actions and liabilities of the Sublessee including but not
limited to: damage to the Premises, non-payment of rent, and any eviction process (In the event
of an eviction the Tenant(s) shall be responsible for all court filing fee(s), representation, and
any other fee(s) associated with removing the Sublessee). The consent by the Landlord to one
sub-let shall not be deemed to be consent to any subsequent subletting.`}}

{{template "paragraph-with-header" dict "Header" "ABANDONMENT" "Text" (printf `If the Tenant(s) vacates or abandons the property for a time-period that is
the minimum set by State law or %s (%d) days, whichever is less, the Landlord shall have the
right to terminate this Agreement immediately and remove all belongings including any
personal property off of the Premises. If the Tenant(s) vacates or abandons the property, the
Landlord shall immediately have the right to terminate this Agreement.` (numberToWord .MinimumAbandonmentDays false) (.MinimumAbandonmentDays))}}

{{template "paragraph-with-header" dict "Header" "ASSIGNMENT" "Text" `Tenant(s) shall not assign this Lease without the prior written consent of the
Landlord. The consent by the Landlord to one assignment shall not be deemed to be consent to
any subsequent assignment.`}}

{{template "paragraph-with-header" dict "Header" "PARKING" "Text" (printf "The Landlord shall provide the Tenant(s) %d Parking Spaces." .Apartment.ParkingSpaces)}}

{{template "simple-paragraph" dict "Text" (printf `The Landlord shall not charge a fee for the %d Parking Spaces. The Parking Space(s) can be
described as: %s provided` .Apartment.ParkingSpaces .Apartment.ParkingSpaceDesc)}}

{{template "paragraph-with-header" dict "Header" "RIGHT OF ENTRY" "Text" `The Landlord shall have the right to enter the Premises during normal
working hours by providing notice in accordance with the minimum State requirement in order
for inspection, make necessary repairs, alterations or improvements, to supply services as
agreed or for any reasonable purpose. The Landlord may exhibit the Premises to prospective
purchasers, mortgagees, or lessees upon reasonable notice.`}}

{{template "paragraph-with-header" dict "Header" "SALE OF PROPERTY" "Text" `If the Premises is sold, the Tenant(s) is to be notified of the new
Owner, and if there is a new Manager, their contact details for repairs and maintenance shall be
forwarded. If the Premises is conveyed to another party, the new owner shall not have the right
to terminate this Agreement and it shall continue under the terms and conditions agreed upon
by the Landlord and Tenant(s).`}}

{{template "paragraph-with-header" dict "Header" "UTILITIES" "Text" "The Landlord agrees to pay for the following utilities and services:"}}

{{template "simple-paragraph" dict "Text" (printf `%s and the Landlord shall also provide Some
great services with all other utilities and services to be the responsibility of the Tenant(s).` (listItems .Apartment.Utilities false))}}

{{template "paragraph-with-header" dict "Header" "MAINTENANCE, REPAIRS, OR ALTERATIONS" "Text" `The Tenant(s) shall, at their own
expense and at all times, maintain the Premises in a clean and sanitary manner, and shall
surrender the same at termination hereof, in as good condition as received, normal wear and
tear excepted. The Tenant(s) may not make any alterations to the leased Premises without the
consent in writing of the Landlord. The Landlord shall be responsible for repairs to the interior
and exterior of the building. If the Premises includes a washer, dryer, freezer, dehumidifier unit
and/or air conditioning unit, the Landlord makes no warranty as to the repair or replacement of
units if one or all shall fail to operate. The Landlord will place fresh batteries in all
battery-operated smoke detectors when the Tenant(s) moves into the Premises. After the initial
placement of the fresh batteries, it is the responsibility of the Tenant(s) to replace batteries
when needed. A monthly “cursory” inspection may be required for all fire extinguishers to
make sure they are fully charged.`}}

{{template "paragraph-with-header" dict "Margin" "20 0 0 0" "Header" "EARLY TERMINATION" "Text" `The Tenant(s) may be allowed to cancel this Agreement under the
following conditions:`}}

{{template "simple-paragraph" dict "Margin" "10 0 0 0" "Text" (printf `The Tenant(s) must provide at least %d days' notice and pay an early termination fee of
$%s (US Dollars) which does not include the rent due for the notice period. During the
notice period of %d days the rent shall be paid in accordance with this Agreement.` .TerminationNoticePeriod .TerminationFee .TerminationNoticePeriod)}}

{{template "paragraph-with-header" dict "Header" "PETS" "Text" "The Tenant(s) shall be allowed to have:"}}

{{template "simple-paragraph" dict "Text" (printf `%s(%d) pets on the Premises consisting of %s, with 
no other types of Pet(s) being allowed on the Premises or common areas, hereinafter known as 
the “Pet(s)”. The Tenant(s) shall be required to pay a pet fee in the amount of $%s for all the Pet(s) which is refundable at the end of the Lease Term only if there is no damage to the
Premises that is caused by the Pet(s). The Tenant(s) is responsible for all damage that any pet
causes, regardless of ownership of said pet and agrees to restore the property to its original
condition at their expense. There shall be no limit on the weight of the pet. pounds (Lb.).` (numberToWord .Apartment.NumberOfAllowedPets  true) (.Apartment.NumberOfAllowedPets) (listItems .Apartment.AllowedPets true) (.PetFee))}}

{{template "paragraph-with-header" dict "Header" "NOISE/WASTE" "Text" `The Tenant(s) agrees not to commit waste on the Premises, maintain, or
permit to be maintained, a nuisance thereon, or use, or permit the Premises to be used, in an
unlawful manner. The Tenant(s) further agrees to abide by any and all local, county, and State
noise ordinances.`}}

{{template "paragraph-with-header" dict "Header" "GUESTS" "Text" `There shall be no other persons living on the Premises other than the Tenant(s) and
any Occupant(s). Guests of the Tenant(s) are allowed for periods not lasting for more than
forty-eight hours unless otherwise approved by the Landlord.
`}}

{{template "paragraph-with-header" dict "Header" "SMOKING POLICY" "Text" `Smoking on the Premises is prohibited on the entire property, including
individual units, common areas, every building and adjoining properties.`}}

{{template "paragraph-with-header" dict "Header" "COMPLIANCE WITH LAW" "Text" `he Tenant(s) agrees that during the term of the Agreement, to
promptly comply with any present and future laws, ordinances, orders, rules, regulations, and
requirements of the Federal, State, County, City, and Municipal government or any of their
departments, bureaus, boards, commissions and officials thereof with respect to the Premises,
or the use or occupancy thereof, whether said compliance shall be ordered or directed to or
against the Tenant(s), the Landlord, or both.`}}

{{template "paragraph-with-header" dict "Header" "DEFAULT" "Text" `If the Tenant(s) fails to comply with any of the financial or material provisions of
this Agreement, or of any present rules and regulations or any that may be hereafter prescribed
by the Landlord, or materially fails to comply with any duties imposed on the Tenant(s) by
statute or State laws, within the time period after delivery of written notice by the Landlord
specifying the non-compliance and indicating the intention of the Landlord to terminate the
Agreement by reason thereof, the Landlord may terminate this Agreement. If the Tenant(s) fails
to pay rent when due and the default continues for the time-period specified in the written
notice thereafter, the Landlord may, at their option, declare the entire balance (compiling all
months applicable to this Agreement) of rent payable hereunder to be immediately due and
payable and may exercise any and all rights and remedies available to the Landlord at law or in
equity and may immediately terminate this Agreement.`}}

{{template "simple-paragraph" dict "Text" `The Tenant(s) will be in default if: (a) Tenant(s) does not pay rent or other amounts that are
owed in accordance with respective State laws; (b) Tenant(s), their guests, or the Occupant(s)
violate this Agreement, rules, or fire, safety, health, or criminal laws, regardless of whether
arrest or conviction occurs; (c) Tenant(s) abandons the Premises; (d) Tenant(s) gives incorrect
or false information in the rental application; (e) Tenant(s), or any Occupant(s) is arrested,
convicted, or given deferred adjudication for a criminal offense involving actual or potential
physical harm to a person, or involving possession, manufacture, or delivery of a controlled
substance, marijuana, or drug paraphernalia under state statute; (f) any illegal drugs or
paraphernalia are found in the Premises or on the person of the Tenant(s), guests, or
Occupant(s) while on the Premises and/or; (g) as otherwise allowed by law`}}

{{template "paragraph-with-header" dict "Header" "MULTIPLE TENANT(S) OR OCCUPANT(S)" "Text" `Each individual that is considered a Tenant(s) is jointly and individually liable for all of this Agreement's obligations, including but
not limited to rent monies. If any Tenant(s), guest, or Occupant(s) violates this Agreement, the
Tenant(s) is considered to have violated this Agreement. Landlord’s requests and notices to the
Tenant(s) or any of the Occupant(s) of legal age constitutes notice to the Tenant(s). Notices and requests from the Tenant(s) or any one of the Occupant(s) (including repair requests and entry
permissions) constitutes notice from the Tenant(s). In eviction suits, the Tenant(s) is considered
the agent of the Premise for the service of process.`}}

{{template "paragraph-with-header" dict "Header" "DISPUTES" "Text" `If a dispute arises during or after the term of this Agreement between the
Landlord and Tenant(s), they shall agree to hold negotiations amongst themselves, in “good
faith”, before any litigation.`}}
{{template "paragraph-with-header" dict "Header" "SEVERABILITY" "Text" `If any provision of this Agreement or the application thereof shall, for any
reason and to any extent, be invalid or unenforceable, neither the remainder of this Agreement
nor the application of the provision to other persons, entities or circumstances shall be affected
thereby, but instead shall be enforced to the maximum extent permitted by law.`}}

{{template "paragraph-with-header" dict "Header" "SURRENDER OF PREMISES" "Text" `The Tenant(s) has surrendered the Premises when (a) the moveout date has passed and no one is living in the Premise within the Landlord’s reasonable
judgment; or (b) Access to the Premise have been turned in to Landlord – whichever comes
first. Upon the expiration of the term hereof, the Tenant(s) shall surrender the Premise in better
or equal condition as it were at the commencement of this Agreement, reasonable use, wear and
tear thereof, and damages by the elements excepted.`}}

{{template "paragraph-with-header" dict "Header" "RETALIATION" "Text" `The Landlord is prohibited from making any type of retaliatory acts against
the Tenant(s) including but not limited to restricting access to the Premises, decreasing or
cancelling services or utilities, failure to repair appliances or fixtures, or any other type of act
that could be considered unjustified.`}}

{{template "paragraph-with-header" dict "Header" "WAIVER" "Text" `A Waiver by the Landlord for a breach of any covenant or duty by the Tenant(s),
under this Agreement is not a waiver for a breach of any other covenant or duty by the
Tenant(s), or of any subsequent breach of the same covenant or duty. No provision of this
Agreement shall be considered waived unless such a waiver shall be expressed in writing as a
formal amendment to this Agreement and executed by the Tenant(s) and Landlord.`}}

{{template "paragraph-with-header" dict "Header" "EQUAL HOUSING" "Text" `If the Tenant(s) possess(es) any mental or physical impairment, the
Landlord shall provide reasonable modifications to the Premises unless the modifications
would be too difficult or expensive for the Landlord to provide. Any impairment of the
Tenant(s) is/are encouraged to be provided and presented to the Landlord in writing in order to
seek the most appropriate route for providing the modifications to the Premises.`}}

{{template "paragraph-with-header" dict "Header" "HAZARDOUS MATERIALS" "Text" `HAZARDOUS MATERIALS: </text-chunk>
<text-chunk font="times" font-size="12">The Tenant(s) agrees to not possess any type of personal
property that could be considered a fire hazard such as a substance having flammable or explosive characteristics on the Premises. Items that are prohibited to be brought into the
Premises, other than for everyday cooking or the need of an appliance, includes but is not
limited to gas (compressed), gasoline, fuel, propane, kerosene, motor oil, fireworks, or any
other related content in the form of a liquid, solid, or gas.`}}

{{template "paragraph-with-header" dict "Header" "WATERBEDS" "Text" "The Tenant(s) is not permitted to furnish the Premises with waterbeds."}}

{{template "paragraph-with-header" dict "Header" "INDEMNIFICATION" "Text" `The Landlord shall not be liable for any damage or injury to the
Tenant(s), or any other person, or to any property, occurring on the Premises, or any part
thereof, or in common areas thereof, and the Tenant(s) agrees to hold the Landlord harmless
from any claims or damages unless caused solely by the Landlord's negligence. It is
recommended that renter's insurance be purchased at the Tenant(s)'s expense.`}}

{{template "paragraph-with-header" dict "Header" "COVENANTS" "Text" `The covenants and conditions herein contained shall apply to and bind the
heirs, legal representatives, and assigns of the parties hereto, and all covenants are to be
construed as conditions of this Agreement.`}}

{{template "paragraph-with-header" dict "Header" "NOTICES" "Text" `Any notice to be sent by the Landlord or the Tenant(s) to each other shall use the
following mailing addresses:`}}

{{template "simple-paragraph" dict "Font" "times-bold" "Text" "Landlord's/Agent's Mailing Address"}}
{{template "simple-paragraph" dict "Text" (printf `%s, ATTN. %s
%s` .Company.Name .Company.LandLord .Company.Address)}}

{{template "simple-paragraph" dict "Font" "times-bold" "Text" "Tenant(s)'s Mailing Address"}}
{{template "simple-paragraph" dict "Text" (printf `%s
%s` (listItems .Tenant.Names true) (.Tenant.MailingAddress))}}

{{template "paragraph-with-header" dict "Header" "AGENT/MANAGER" "Text" (printf `The Landlord authorizes the following to act on their behalf in regards
to the Premises for any repair, maintenance, or compliant other than a breach of this
Agreement: The The management company known as %s of %s that can be contacted at the following Phone
Number %s and can be E-Mailed at %s.` .Manager.Company .Manager.Address .Manager.Phone .Manager.Email)}}

{{template "paragraph-with-header" dict "Header" "PREMISES DEEMED UNINHABITABLE" "Text" `If the Property is deemed uninhabitable due to
damage beyond reasonable repair the Tenant(s) will be able to terminate this Agreement by
written notice to the Landlord. If said damage was due to the negligence of the Tenant(s), the
Tenant(s) shall be liable to the Landlord for all repairs and for the loss of income due to
restoring the Premises back to a livable condition in addition to any other losses that can be
proved by the Landlord.`}}

{{template "paragraph-with-header" dict "Header" "SERVICEMEMBERS CIVIL RELIEF ACT" "Text" (printf `In the event the Tenant(s) is or hereafter
becomes, a member of the United States Armed Forces on extended active duty and hereafter
the Tenant(s) receives permanent change of station (PCS) orders to depart from the area where
the Premises are located, or is relieved from active duty, retires or separates from the military,
is ordered into military housing, or receives deployment orders, then in any of these events, the
Tenant may terminate this lease upon giving %s (%d) days written notice to the Landlord.
The Tenant shall also provide to the Landlord a copy of the official orders or a letter signed by
the Tenant’s commanding officer, reflecting the change which warrants termination under this
clause. The Tenant will pay prorated rent for any days which he/she occupies the dwelling past the beginning of the rental period.` (numberToWord .LeaseTerminationOfServiceMembers false) (.LeaseTerminationOfServiceMembers))}}

{{template "simple-paragraph" dict "Text" "The damage/security deposit will be promptly returned to Tenant, provided there are no damages to the Premises"}}

{{template "paragraph-with-header" dict "Header" "LEAD PAINT" "Text" (printf `The Premises was not constructed before %s and therefore does not contain
leadbased paint.` .Apartment.ConstructedBefore)}}

{{template "paragraph-with-header" dict "Header" "GOVERNING LAW" "Text" (printf `This Agreement is to be governed under the laws located in the State of
%s` .Company.Location)}}

{{template "paragraph-with-header" dict "Header" "ADDITIONAL TERMS AND CONDITIONS" "Text" `In addition to the above stated terms and
conditions of this Agreement, the Landlord and Tenant agree to the following: Additional
Terms are to be specified: Term 1, Term 2, Term 3`}}

{{template "paragraph-with-header" dict "Margin" "18 0 20 0" "Header" "ENTIRE AGREEMENT" "Text" (printf `This Agreement contains all the terms agreed to by the parties
relating to its subject matter including any attachments or addendums. This Agreement replaces
all previous discussions, understandings, and oral agreements. The Landlord and Tenant(s)
agree to the terms and conditions and shall be bound until the end of the Lease Term.

The parties have agreed and executed this agreement on %s` (formatTime .Date "December 9 2006"))}}

{{/* render on a new page */}}
<page-break></page-break>
{{template "simple-paragraph" dict "Margin" "18 0 10 0" "Font" "times-bold" "Text" "LANDLORD(S) SIGNATURE"}}
<division margin="20 60 0 0">
   {{template "form-sig" dict "Margin" "0 0 0 110" "Text" "Landlord’s Signature"}}
</division>
{{template "simple-paragraph" dict "Margin" "5 0 0 0" "Text" (printf `%s as President of %s` .Company.LandLord .Company.Name)}}
{{template "simple-paragraph" dict "Margin" "18 0 20 0" "Font" "times-bold" "Text" "TENANT(S) SIGNATURE"}}
<division margin="20 60 0 0">
   {{template "form-sig" dict "Margin" "0 0 0 100" "Text" "Tenant’s Signature"}}
</division>
<division margin="20 60 20 0">
   {{template "form-sig" dict "Margin" "0 0 0 100" "Text" "Tenant’s Signature"}}
</division>

<page-break></page-break>
<paragraph margin="18 0 10 0" text-align = "center">
   <text-chunk font="times-bold" font-size="21.5">Security Deposit Receipt</text-chunk>
</paragraph>

<division margin="10 0 0 0">
   <paragraph line-height="2.5">
      <text-chunk font="times">Date:__________________________________________{{printf "\n"}}</text-chunk>
      <text-chunk font="times">Dear ___________________________________________[Tenant(s)],{{printf "\n"}}</text-chunk>
      <text-chunk font="times">The Landlord shall hold the Security Deposit in a separate account at a bank{{printf "\n"}}</text-chunk>
      <text-chunk font="times">located at __________________________________________[Street Address] in{{printf "\n"}}</text-chunk>
      <text-chunk font="times">the City of __________________________________________ , State of ___________________{{printf "\n"}}</text-chunk>
      <text-chunk font="times">The Security Deposit in the amount of $ _____________________ (US Dollars) has been deposited in{{printf "\n"}}</text-chunk>
      <text-chunk font="times">___________________ [Bank Name] with the Account Number of _________________ for the full{{printf "\n"}}</text-chunk>
      <text-chunk font="times">performance of the Lease executed on the _____ day of _______________ , 20 ___.{{printf "\n"}}</text-chunk>
      <text-chunk font="times">Sincerely,{{printf "\n"}}</text-chunk>
   </paragraph>
</division>
<division margin="20 0 30 0">
   {{template "form-sig" dict "Margin" "0 200 0 110" "Text" "Landlord’s Signature"}}
</division>
<page-break></page-break>
<division margin="50 0 30 0">
   <paragraph text-align = "center">
      <text-chunk font="times" font-size="12">AMOUNT ($) DUE AT SIGNING</text-chunk>
   </paragraph>
</division>
<paragraph margin="18 0" line-height="2.3">
   <text-chunk font="times-bold" font-size="12">Security Deposit: </text-chunk> 
   <text-chunk font="times" font-size="12">${{.SecurityDeposit}}{{printf "\n"}}</text-chunk>
   <text-chunk font="times-bold" font-size="12">First (1st) Month's Rent: </text-chunk>
   <text-chunk font="times" font-size="12">${{.SecurityDeposit}}{{printf "\n"}}</text-chunk>
   <text-chunk font="times-bold" font-size="12">Pet Fee(s): </text-chunk> 
   <text-chunk font="times" font-size="12"> ${{.PetFee}} for all the Pet(s)</text-chunk>
</paragraph>
<page-break></page-break>
<paragraph text-align="center">
   <text-chunk font="times-bold" font-size="18.5">Move-in Checklist</text-chunk>
</paragraph>
<paragraph>
   <text-chunk font="times" font-size="11">Property Address: {{.Apartment.Address}}{{printf "\n"}}Unit Size: {{.Apartment.UnitSize}} bedroom(s){{printf "\n"}}</text-chunk>
   <text-chunk font="times" font-size="11">Move-in Inspection Date: ___________________ Move-out Inspection Date: _________________</text-chunk>                             
</paragraph>

{{template "simple-paragraph" dict "Margin" "0" "FontSize" 11 "Text" `
Write the condition of the space along with any specific damage or repairs needed. Be sure to write
any repair needed such as paint chipping, wall damage, or any lessened area that could be considered
maintenance needed at the end of the lease, and therefore, be deducted at the end of the Lease Term.`}}

{{template "simple-paragraph" dict "FontSize" 18.5 "Margin" "15 0 0 0" "Font" "times-bold" "Text" "Living Room"}}
{{template "check-list-table" dict "Items" .MoveInCheckList.LivingRoom}}
{{template "simple-paragraph" dict "FontSize" 18.5 "Margin" "10 0 0 0" "Font" "times-bold" "Text" "Dining Room"}}
{{template "check-list-table" dict "Items" .MoveInCheckList.DinningRoom}}
{{template "simple-paragraph" dict "FontSize" 18.5 "Margin" "10 0 0 0" "Font" "times-bold" "Text" "Kitchen Area"}}
{{template "check-list-table" dict "Items" .MoveInCheckList.Kitchen}}
{{template "simple-paragraph" dict "FontSize" 18.5 "Margin" "10 0 0 0" "Font" "times-bold" "Text" "Bedroom(s)"}}
{{template "check-list-table" dict "Items" .MoveInCheckList.Bathroom}}
{{template "simple-paragraph" dict "FontSize" 18.5 "Margin" "10 0 0 0" "Font" "times-bold" "Text" "Other"}}
{{template "check-list-table" dict "Items" .MoveInCheckList.Other}}

{{template "simple-paragraph" dict "Text" `I, a Tenant on this Lease, have sufficiently inspected the Premises and confirm above-stated 
information. (only 1 Tenant required)`}}
<division margin="18 60 0 0">
   {{template "form-sig" dict "Margin" "0 0 0 100" "Text" "Tenant’s Signature"}}
</division>

{{template "simple-paragraph" dict "Text" `I, the Landlord on this Lease, have sufficiently inspected the Premises and confirm 
above-statedinformation`}}

<division margin="18 60 0 0">
   {{template "form-sig" dict "Margin" "0 0 0 110" "Text" "Landlord’s Signature"}}
</division>