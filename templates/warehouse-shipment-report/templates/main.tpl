{{define "shipment"}}
    <list-item>
        <list-marker font-size="25" font="helvetica-bold">• </list-marker>

        <paragraph>
            <text-chunk color="#3C1C00" font-size="25" font="helvetica-bold">{{ printf "Shipment %s" .Shipment.Time }}</text-chunk>
        </paragraph>
    </list-item>

    <list-item>
        <list-marker> </list-marker>

        <list indent="0">
            {{range $idx, $order := .Shipment.Orders}}
                {{template "order" (dict "OrderIndex" (printf "%d.%d" $.ShipmentIndex (add $idx 1)) "Order" $order ) }}
            {{end}}
        </list>
    </list-item>
{{end}}

{{define "order"}}
    <list-item>
        <list-marker font-size="15">{{ .OrderIndex }}. </list-marker>
        
        <paragraph>
            <text-chunk font-size="15">{{ printf "Order %s, %s" .Order.Time .Order.PIC }}</text-chunk>
        </paragraph>
    </list-item>

    <list-item>
        <list-marker> </list-marker>

        <list indent="0">
            {{range $idx, $product := .Order.Products}}
                {{template "product" (dict "ProductIndex" (printf "%s.%d" $.OrderIndex (add $idx 1)) "Product" $product ) }}
            {{end}}
        </list>
    </list-item>
{{end}}

{{define "product"}}
    <list-item>
        <list-marker>{{ .ProductIndex }}. </list-marker>

        <list indent="-5">
            <list-marker>- </list-marker>

            <list-item>
                <list-marker> </list-marker>
                
                <image src="{{ createBarcode .Product.Barcode }}"></image>
            </list-item>

            <list-item>
                <paragraph>
                    <text-chunk>Product Code: {{ .Product.Code }}</text-chunk>
                </paragraph>
            </list-item>

            <list-item>
                <paragraph>
                    <text-chunk>Product Name: {{ html .Product.Name }}</text-chunk>
                </paragraph>
            </list-item>
        </list>
    </list-item>
{{end}}

{{range $idx, $shipment := .}}
    <list>
        {{template "shipment" (dict "ShipmentIndex" (add $idx 1) "Shipment" $shipment ) }}
    </list>
{{end}}